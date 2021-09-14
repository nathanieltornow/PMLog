#include <cstdlib>
#include <iostream>
#include <ctime>
#include <utility>
#include <cstdint>
#include <libpmemobj++/container/concurrent_hash_map.hpp>
#include <libpmemobj++/p.hpp>
#include <libpmemobj++/persistent_ptr.hpp>
#include <libpmemobj++/make_persistent_array_atomic.hpp>
#include <libpmemobj++/transaction.hpp>
#include <libpmemobj++/pool.hpp>
#include "PMLog.hpp"

#define KEY_SIZE 16

using namespace pmem::obj;

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

void *cppStartUp() {
	//auto pop = pmem::obj::pool<root>::create("log", "example", 33554432);
	pool<root> pop;
	
	try {
		if (pool<root>::check("log", "example") == 1)
			pop = pool<root>::open("log", "example");
		else {
			std::cerr << "Memory pool is corrupted. Exiting...\n";
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

void cppFinalize(persistent_ptr<cppPMLog> cppLog) {
	pool<root> pop = cppLog->pop.get_rw();
	
	
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
			persistent_ptr<PString> record = acc->second;
			if (!(gsnPptr->insert(GSNmap::value_type(gsn, record))))
				return -2;
			
			this->highest_gsn = gsn;
			this->pop.get_rw().persist(this->highest_gsn);
			
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
	if (gsn > this->highest_gsn.get_ro())
		return 0;
	
	uint64_t next_gsn = gsn + 1;
	
	try {
		GSNmap::accessor acc;
		if (gsnPptr->find(acc, gsn))
			strcpy(storage, acc->second->data());
			
		uint64_t next_gsn = gsn + 1;
		while (!(gsnPptr->find(acc, next_gsn))) {
			if (next_gsn > this->highest_gsn)
				return 0;
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
		for (it = gsnPptr->begin(); it != gsnPptr->end(); it++) {
			if (it->first < gsn) {
				gsnPptr->erase(it->first);
			}
		}
		
		return 0;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		this->pop.get_rw().close();
		exit(-1);
	}	
}

using cppPMLog = struct cppPMLog;
