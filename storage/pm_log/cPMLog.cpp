#include "PMLog.hpp"
#include "PMLog.h"

void PMLogFree(PMLog log) {
    cppPMLog *cppLog = (cppPMLog *) log;

    delete cppLog;
}

PMLog startUp() {
	return (PMLog) cppStartUp();
}

int Append(PMLog log, char* record, uint64_t lsn) {
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return cppLog->Append(record, lsn);
}

int Commit(PMLog log, uint64_t lsn, uint32_t color, uint64_t gsn){
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return cppLog->Commit(lsn, color, gsn);
}

const char* Read(PMLog log, uint32_t color, uint64_t gsn) {
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return (cppLog->Read(color, gsn).c_str());
}

int Trim(PMLog log, uint32_t color, uint64_t gsn) {
    persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

    return cppLog->Trim(color, gsn);
}