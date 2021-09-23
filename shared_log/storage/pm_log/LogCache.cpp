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
		if (this->logCachePtr->size() == MAX_CACHE_SIZE) {
			this->logCachePtr->erase(this->lowest_gsn);
			uint64_t tmp = this->next_gsn(this->lowest_gsn);
			
			if (tmp != gsn)
				this->lowest_gsn = tmp;
		}
		
        if (!(this->logCachePtr->insert(std::make_pair(gsn, record))))
			return -1;		
		
		this->highest_gsn = gsn;
		
        return 0;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}

uint64_t LogCache::Read(uint64_t gsn, char *storage) {
	GSNmapCache::accessor acc;

    try {
		uint64_t ret = this->next_gsn(gsn);
		
        if (!(this->logCachePtr->find(acc, gsn)))
            return ret;
        
        strcpy(storage, acc->second.c_str());

        return ret;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}