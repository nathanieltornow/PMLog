#include "Log.hpp"
#include "Log.h"

Log LogNew() {
    return (Log *) new CppLog();
}

void LogFree(Log log) {
    CppLog *cppLog = (CppLog *) log;

    delete cppLog;
}

int cAppend(Log log, const char* record, uint64_t lsn) {	
    CppLog *cppLog = (CppLog *) log;

    return cppLog->Append(record, lsn);
}

int cCommit(Log log, uint64_t lsn, uint64_t gsn){
    CppLog *cppLog = (CppLog *) log;

    return cppLog->Commit(lsn, gsn);
}

uint64_t cRead(Log log, uint64_t gsn, char* storage) { 
    CppLog *cppLog = (CppLog *) log;

    return cppLog->Read(gsn, storage);
}

int cTrim(Log log, uint64_t gsn) {
   CppLog *cppLog = (CppLog *) log;

    return cppLog->Trim(gsn);
}