#include <cstdlib>
#include <iostream>
#include <ctime>
#include <utility>
#include <cstdint>
#include <libpmemobj++/container/concurrent_hash_map.hpp>
#include <libpmemobj++/p.hpp>
#include <libpmemobj++/persistent_ptr.hpp>
#include <libpmemobj++/pool.hpp>
#include "PMLog.hpp"

#define KEY_SIZE 16

using namespace pmem::obj;

struct root {
    persistent_ptr<cppPMLog> pmLog;
};

void *cppStartUp() {
	//auto pop = pmem::obj::pool<root>::create("log", "example", 33554432); 
	auto pop = pool<root>::open("log", "example");
    auto r = pop.root(); 
	
	if (r->pmLog == nullptr) {
		pmem::obj::transaction::run(pop, [&] {
			r->pmLog = make_persistent<cppPMLog>(pop);
		});
	}
	else {
		r->pmLog->restartMaps();
	}
	
	return (void *) (r->pmLog.get());
}

const char* PersistentString::data() const{ 
	return this->array; 
}

PersistentString::PersistentString(const char* s) {	
	pmemobj_tx_add_range_direct(this->array, KEY_SIZE);
	strcpy(this->array, s);
}
	
cppPMLog::cppPMLog(pool_base pop) {
	this->pop = pop;
	this->highest_gsn = 0;

	pmem::obj::transaction::run(pop, [&] {
		lsnPptr = make_persistent<LSNmap>();
		gsnPptr = make_persistent<GSNmap>();
	});
}

cppPMLog::~cppPMLog() {

}

void cppPMLog::restartMaps() {		
	this->lsnPptr->defragment();
	this->gsnPptr->defragment();
	this->lsnPptr->runtime_initialize();
	this->gsnPptr->runtime_initialize();
}

int cppPMLog::Append(const char* record, uint64_t lsn) {	
	persistent_ptr<PString> s; 
	
	transaction::run(this->pop, [&] {
		s = make_persistent<PString>(record);
	});
	this->lsnPptr->insert(LSNmap::value_type(lsn, s));
	
	return 0;
}

int cppPMLog::Commit(uint64_t lsn, uint64_t gsn) {
	LSNmap::accessor acc;
	
	bool res = this->lsnPptr->find(acc, lsn);

	if (res) {
		persistent_ptr<PString> record = acc->second;
		gsnPptr->insert(GSNmap::value_type(gsn, record));
		this->highest_gsn = gsn;
		
		return 0;
	}
	else {
		std::cout << "Record with LSN " << lsn << "not found" << std::endl;
	
		return -1;
	}	
}

uint64_t cppPMLog::Read(uint64_t gsn, char* storage) {
	if (gsn > this->highest_gsn)
		return 0;
	
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

int cppPMLog::Trim(uint64_t gsn) {
	GSNmap::iterator it;
	
	for (it = gsnPptr->begin(); it != gsnPptr->end(); it++) {
		if (it->first < gsn) {
			gsnPptr->erase(it->first);
		}
	}
	
	return 0;
}

using cppPMLog = struct cppPMLog;
	
/*
struct root {
    persistent_ptr<cppPMLog> pmLog;
};

int main() {
	auto pop = pmem::obj::pool<root>::create("log", "example", 33554432); 
	//auto pop = pool<root>::open("log", "example");
    auto r = pop.root(); 
	
	if (r->pmLog == nullptr) {
		pmem::obj::transaction::run(pop, [&] {
			r->pmLog = make_persistent<cppPMLog>(pop);
		});
	}
	else {
		pop.root()->pmLog->restartMaps();
	}
	
	auto pmLog = pop.root()->pmLog;
    auto &LSNMap = pmLog->lsnPptr;
    auto &ColorMap = pmLog->colorPptr;	
	
	for (int i = 0; i < 1000; i++) {
		int lsn = Append(pmLog, std::string("manos" + std::to_string(i)), i);
		Commit(pmLog, i, 0, i); 
	}
	
	Trim(pmLog, 0, 310);
	
	for (int i = 300; i < 330; i++) {
		std::cout << Red(pmLog, 0, i) << std::endl;
	}
	
	Trim(pmLog, 0, 500);
	
	for (int i = 1000; i < 1100; i++) {
		int lsn = Append(pmLog, std::string("manos" + std::to_string(i)), i);
		Commit(pmLog, i, 0, i); 
	}
	
	Trim(pmLog, 0, 1000);
	
	for (int i = 980; i < 1100; i++) {
		std::cout << Read(pmLog, 0, i) << std::endl;
	}
	
	return 0;
}
*/