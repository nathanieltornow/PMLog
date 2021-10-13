#include <cstdlib>
#include <iostream>
#include <fstream>
#include <ctime>
#include <utility>
#include <cstdint>
#include <libpmemobj++/container/concurrent_hash_map.hpp>
#include <libpmemobj++/p.hpp>
#include <libpmemobj++/persistent_ptr.hpp>
#include <libpmemobj++/make_persistent_array_atomic.hpp>
#include <libpmemobj++/transaction.hpp>
#include <libpmemobj++/pool.hpp>
#include <thread>  
#include "PMLog.hpp"
#include "LogCache.hpp"

static size_t CACHE_SEGMENT_SIZE = 0;
extern size_t MAX_CACHE_SIZE;

using namespace pmem::obj;

static LogCache logCache;

struct root {
    persistent_ptr<cppPMLog> pmLog;
};

const char* PersistentString::data() const{ 
	return this->array; 
}

PersistentString::PersistentString(const char* s) {	
	if (pmemobj_tx_add_range_direct(this->array, KEY_SIZE) != 0) {
		std::cerr << "Storage of string " << s << "failed.\n";
	}
	strcpy(this->array, s);
}

void setup(std::string &s1, std::string &s2) {
	std::ifstream logFile;
	logFile.open("configuration.txt", std::ifstream::in);

	if (logFile.is_open()) {
		getline(logFile, s1);
		getline(logFile, s2);
		
		std::string tmp;
		getline(logFile, tmp);
		CACHE_SEGMENT_SIZE = std::stoi(tmp);
		
		getline(logFile, tmp);
		MAX_CACHE_SIZE = std::stoi(tmp);
		
		logFile.close();
	}
}

void *cppStartUp() {
	pool<root> pop;
	std::string s1, s2;
	
	setup(s1, s2);
	
	s1.pop_back();
	s2.pop_back();

	try {
		if (pool<root>::check(s1, s2) == 1)
			pop = pool<root>::open(s1, s2);
		else {
			std::cerr << "Memory pool " << s1 << " with layout " << s2<< " is corrupted or does not exist. Exiting...\n" << std::flush;
			exit(-1);
		}
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}
    auto r = pop.root(); 
	
	if (r->pmLog == nullptr) {		
		try {
			pmem::obj::transaction::run(pop, [&] {
				r->pmLog = make_persistent<cppPMLog>(pop);
			});
		}
		catch (const std::runtime_error &e){
			std::cerr << e.what();
			exit(-1);
		}
	}
	else {
		r->pmLog->restartMaps();
	}
	
	return (void *) (r->pmLog.get());
}

void cppPMLog::shutdown() {
	this->lsnPptr->defragment();
	this->gsnPptr->defragment();
}

void cppFinalize(persistent_ptr<cppPMLog> cppLog) {
	pool<root> pop = cppLog->pop.get_rw();
	cppLog->shutdown();
	
	pmem::obj::transaction::run(pop, [&] {
		delete_persistent<cppPMLog>(cppLog);
		pop.root()->pmLog = nullptr;
	});	

	pop.close();
}

cppPMLog::cppPMLog(pool<root> pop) {
	try {
		pmem::obj::transaction::run(pop, [&] {
			this->pop = pop;
			this->highest_gsn = 0;
			this->lowest_gsn = 0;
			lsnPptr = make_persistent<LSNmap>();
			gsnPptr = make_persistent<GSNmap>();
		});
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		pop.close();
		exit(-1);
	}
}

cppPMLog::~cppPMLog() {
	try {
		this->lsnPptr->clear();
		this->gsnPptr->clear();
		pmem::obj::transaction::run(this->pop.get_rw(), [&] {
			delete_persistent<LSNmap>(this->lsnPptr);
			delete_persistent<GSNmap>(this->gsnPptr);
		});		
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}
}

void cppPMLog::restartMaps() {
	try {
		this->lsnPptr->defragment();
		this->gsnPptr->defragment();
		this->lsnPptr->runtime_initialize();
		this->gsnPptr->runtime_initialize();
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		this->pop.get_rw().close();
		exit(-1);
	}
}

void cppPMLog::cacheRecords(cppPMLog *log, uint64_t gsn) {
	int records = 0;
	uint64_t curr_gsn = gsn;
	GSNmap::accessor acc;
	
	try {
		while (records < CACHE_SEGMENT_SIZE) {
			if (log->gsnPptr->find(acc, curr_gsn))
				logCache.Append(acc->second->data(), curr_gsn);
				records++;
			curr_gsn++;
		}	
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		log->pop.get_rw().close();
		exit(-1);
	}	
}

int cppPMLog::Append(const char* record, uint64_t lsn) {
	persistent_ptr<PString> s;
	bool res;	
	pool<root> p = this->pop.get_rw();
	
	try {
		pmem::obj::transaction::run(this->pop.get_rw(), [&] {
			s = make_persistent<PString>(record);
		});
		res = this->lsnPptr->insert(LSNmap::value_type(lsn, s));
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		this->pop.get_rw().close();
		exit(-1);
	}	
	if (res)	
		return 0;
	else
		return -1;
}

int cppPMLog::Commit(uint64_t lsn, uint64_t gsn) {
	LSNmap::accessor acc;
	bool res;
	try {
		res = this->lsnPptr->find(acc, lsn);

		if (res) {
			uint64_t tmp_highest_gsn = this->highest_gsn.get_ro();
			
			persistent_ptr<PString> record = acc->second;
			if (!(this->gsnPptr->insert(GSNmap::value_type(gsn, record))))
				return -2;
			
			if ( __atomic_compare_exchange_n(&(this->highest_gsn.get_rw()), &tmp_highest_gsn, gsn, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST))
				this->pop.get_rw().persist(this->highest_gsn);
			
			acc.release();
			this->lsnPptr->erase(lsn);
			logCache.Append(record->data(), gsn); 
			
			return 0;
		}
		else
			return -1;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		this->pop.get_rw().close();
		exit(-1);
	}
}

uint64_t cppPMLog::Read(uint64_t gsn, char* storage) {
	if (gsn > this->highest_gsn.get_ro() || gsn < this->lowest_gsn.get_ro())
		return 0;
	
	uint64_t next_gsn = logCache.Read(gsn, storage);
	
	if (next_gsn != gsn)
		return next_gsn;
	else {
		std::thread tmp(&cppPMLog::cacheRecords, this, this, gsn);
		tmp.detach();
	}
	
	try {
		GSNmap::accessor acc;
		if (this->gsnPptr->find(acc, gsn))
			strcpy(storage, acc->second->data());
		
		next_gsn = gsn + 1;
		while (!(this->gsnPptr->find(acc, next_gsn))) {
			if (next_gsn > this->highest_gsn)
				return gsn;
			next_gsn++;
		}
		return next_gsn;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		this->pop.get_rw().close();
		exit(-1);
	}
}

int cppPMLog::Trim(uint64_t gsn) {
	GSNmap::iterator it;
	
	try {
		for (it = this->gsnPptr->begin(); it != this->gsnPptr->end(); it++) {
			if (it->first < gsn) {
				this->gsnPptr->erase(it->first);
				logCache.Erase(it->first);
			}
		}
		
		uint64_t tmp_highest_gsn = this->highest_gsn.get_ro();
		uint64_t tmp_lowest_gsn = this->lowest_gsn.get_ro();
		
		if (gsn == tmp_highest_gsn) {
			if( __atomic_compare_exchange_n(&(this->highest_gsn.get_rw()), &tmp_highest_gsn, 0, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST))
				this->pop.get_rw().persist(this->highest_gsn);
		}
		
		if( __atomic_compare_exchange_n(&(this->lowest_gsn.get_rw()), &tmp_lowest_gsn, gsn, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST))
			this->pop.get_rw().persist(this->lowest_gsn);		
		
		return 0;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		this->pop.get_rw().close();
		exit(-1);
	}	
}

using cppPMLog = struct cppPMLog;