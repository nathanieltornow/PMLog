#include "Log.hpp"
#include <iostream>

CppLog::CppLog() {
    this->lsnPtr = new LSNmap();
    this->gsnPtr = new GSNmap();
}

CppLog::~CppLog() {
    this->lsnPtr->clear();
    this->gsnPtr->clear();
    delete this->lsnPtr;
    delete this->gsnPtr;
}

int CppLog::Append(std::string record, uint64_t lsn) {
    LSNmap::accessor acc;

    try {
        if (this->lsnPtr->find(acc, lsn))
            return -1;
        
        this->lsnPtr->insert(std::make_pair(lsn, record));

        return 0;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}
}

int CppLog::Commit(uint64_t lsn, uint64_t gsn){
    GSNmap::accessor acc;

    try {
        if (this->gsnPtr->find(acc, gsn))
            return -1;

        LSNmap::accessor acc;

        if (!(this->lsnPtr->find(acc, lsn)))
            return -2;

        this->gsnPtr->insert(std::make_pair(gsn, acc->second));
        this->highest_gsn = gsn;

        return 0;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}
}

uint64_t CppLog::Read(uint64_t gsn, char *storage) {
    GSNmap::accessor acc;

    try {
        if (gsn > this->highest_gsn || !(this->gsnPtr->find(acc, gsn)))
            return 0;
        
        strcpy(storage, acc->second.c_str());

        uint64_t next_gsn = gsn + 1;
        while (!(this->gsnPtr->find(acc, next_gsn))) {
            if (next_gsn > this->highest_gsn)
                return 0;
            next_gsn++;
        }
        return next_gsn;
    }
    catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}

int CppLog::Trim(uint64_t gsn){
    GSNmap::iterator it;
	
	try {
		for (it = gsnPtr->begin(); it != gsnPtr->end(); it++) {
			if (it->first < gsn) {
				gsnPtr->erase(it->first);
			}
		}
		
		return 0;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}