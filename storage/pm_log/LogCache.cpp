#include "LogCache.hpp"
#include <iostream>

size_t MAX_CACHE_SIZE = 0;

LogCache::LogCache() {
	this->logCachePtr = new GSNmapCache();
	this->highest_gsn = 0;
}

LogCache::~LogCache() {
	this->logCachePtr->clear();
    delete this->logCachePtr;
}

uint64_t LogCache::next_gsn(uint64_t gsn) {
	GSNmapCache::accessor acc;
	uint64_t next_gsn = gsn + 1;
	
	if (next_gsn > this->highest_gsn)
		return gsn;
	
	while (!(this->logCachePtr->find(acc, next_gsn)))
		next_gsn++;
	
	return next_gsn;
}

int LogCache::Append(std::string record, uint64_t gsn) {
	try {		
		while (this->logCachePtr->size() >= MAX_CACHE_SIZE) {
			
			if (this->logCachePtr->erase(this->lowest_gsn)) {
				uint64_t tmp = this->next_gsn(this->lowest_gsn);
				uint64_t tmp_lowest_gsn = this->next_gsn(this->lowest_gsn);
			
				__atomic_compare_exchange_n(&(this->lowest_gsn), &tmp_lowest_gsn, tmp, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST);
			}
		}
		
        if (!(this->logCachePtr->insert(std::make_pair(gsn, record))))
			return -1;		
		
		uint64_t tmp_highest_gsn = this->highest_gsn;
		__atomic_compare_exchange_n(&(this->highest_gsn), &tmp_highest_gsn, gsn, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST);
		
        return 0;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}

const char *LogCache::Read(uint64_t gsn, uint64_t *next_gsn) {
	GSNmapCache::accessor acc;

    try {
		uint64_t ret = this->next_gsn(gsn);
		
        if (!(this->logCachePtr->find(acc, gsn)))
            return "";
      
		*next_gsn = ret;
		
        return acc->second.c_str();
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}

int LogCache::Erase(uint64_t gsn) {
	 try {
        if (!(this->logCachePtr->erase(gsn)))
            return 0;
		else
			return 1;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}