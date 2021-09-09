#include "PMLog.hpp"
#include "PMLog.h"

void PMLogFree(PMLog log) {
    cppPMLog *cppLog = (cppPMLog *) log;

    delete cppLog;
}

PMLog startUp() {
	return (PMLog) cppStartUp();
}

int cAppend(PMLog log, const char* record, uint64_t lsn) {
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);
	
    return cppLog->Append(record, lsn);
}

int cCommit(PMLog log, uint64_t lsn, uint64_t gsn){
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return cppLog->Commit(lsn, gsn);
}

uint64_t cRead(PMLog log, uint64_t gsn, char* storage) {
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return cppLog->Read(gsn, storage);
}

int cTrim(PMLog log, uint64_t gsn) {
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return cppLog->Trim(gsn);
}