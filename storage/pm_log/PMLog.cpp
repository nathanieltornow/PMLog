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

std::string PersistentString::data() const { 
	return std::string(this->array); 
}

PersistentString::PersistentString(std::string* s) {
	unsigned long length = s->length();
	
	if (length < KEY_SIZE - 1) {
		pmemobj_tx_add_range_direct(this->array, KEY_SIZE);
		strcpy(this->array, s->c_str());
	}
}


using PString = struct PersistentString;
using LSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;
using GSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;
using ColorMap = concurrent_hash_map<p<uint32_t>, persistent_ptr<GSNmap>>;
	
cppPMLog::cppPMLog(pool_base pop) {
	this->pop = pop;

	pmem::obj::transaction::run(pop, [&] {
		lsnPptr = make_persistent<LSNmap>();
		colorPptr = make_persistent<ColorMap>();
	});
}

cppPMLog::~cppPMLog() {

}

void cppPMLog::restartMaps() {		
	this->lsnPptr->defragment();
	this->colorPptr->defragment();
	this->lsnPptr->runtime_initialize();
	this->colorPptr->runtime_initialize();
	
	ColorMap::iterator it;
	
	for (it = this->colorPptr->begin(); it != this->colorPptr->end(); it++) {
		it->second->defragment();
		it->second->runtime_initialize();			
	}
}

int cppPMLog::Append(std::string record, uint64_t lsn) {		
	persistent_ptr<PString> s;
	transaction::run(this->pop, [&] {
		s = make_persistent<PString>(&record);
	});
	this->lsnPptr->insert(LSNmap::value_type(lsn, s));
	
	LSNmap::accessor acc;
	this->lsnPptr->find(acc, lsn);
	
	persistent_ptr<PString> r = acc->second;

	return 0;
}

int cppPMLog::Commit(uint64_t lsn, uint32_t color, uint64_t gsn) {
	LSNmap::accessor acc;
	
	bool res = this->lsnPptr->find(acc, lsn);

	if (res) {
		persistent_ptr<PString> record = acc->second;
		ColorMap::accessor cacc;
		if (!(this->colorPptr->find(cacc, color))) {
			persistent_ptr<GSNmap> gsnMap;
			transaction::run(this->pop, [&] {
				gsnMap = make_persistent<GSNmap>();
			});
			this->colorPptr->insert(ColorMap::value_type(color, gsnMap));
		}
		this->colorPptr->find(cacc, color);
		persistent_ptr<GSNmap> gsnMap = cacc->second;
		gsnMap->insert(GSNmap::value_type(gsn, record));

		return 0;
	}
	else {
		std::cout << "Record with LSN " << lsn << "not found" << std::endl;
	
		return -1;
	}	
}

std::string cppPMLog::Read(uint32_t color, uint64_t gsn) {
	ColorMap::accessor cacc;

	bool res = this->colorPptr->find(cacc, color);

	if (res) {
		persistent_ptr<GSNmap> gsnMap = cacc->second;
		GSNmap::accessor acc;
		if (gsnMap->find(acc, gsn)) {
			return acc->second->data();
		}
		else {
			std::cout << "Record with GSN " << gsn << " not found " << std::endl;
			return "";
		}
	}
	else {
		std::cout << "Log with color " << color << " not found " << std::endl;
	
		return "";
	}	
}

int cppPMLog::Trim(uint32_t color, uint64_t gsn) {
	ColorMap::accessor cacc;
	
	bool res = this->colorPptr->find(cacc, color);
	
	if (res) {
		persistent_ptr<GSNmap> gsnMap = cacc->second;
		GSNmap::iterator it;
		
		for (it = gsnMap->begin(); it != gsnMap->end(); it++) {
			if (it->first < gsn) {
				gsnMap->erase(it->first);
			}
		}
		
		return 0;
	}
	else {
		std::cout << "Cannot find log with Color " << color << std::endl;
		return -1;
	}	
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