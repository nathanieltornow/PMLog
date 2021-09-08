
#include <stdint.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif
    typedef void* PMLog;
	PMLog startUp();
    void PMLogFree(PMLog log);
    int Append(PMLog log, char* record, uint64_t lsn);
    int Commit(PMLog log,uint64_t lsn, uint32_t color, uint64_t gsn) ;
    const char* Read(PMLog log, uint32_t color, uint64_t gsn);
    int Trim(PMLog log, uint32_t color, uint64_t gsn);
#ifdef __cplusplus
}
#endif