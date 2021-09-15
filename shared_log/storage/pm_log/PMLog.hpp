#include <libpmemobj++/persistent_ptr.hpp>
#include <libpmemobj++/container/concurrent_hash_map.hpp>
#include <cstdint>
#include <string>

#define KEY_SIZE 16
using namespace pmem::obj;

struct root;

class PersistentString {
private:
	char array[KEY_SIZE];
public:
	const char* data() const;	
	PersistentString(const char* s);
};

using PString = struct PersistentString;
using LSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;
using GSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;

class cppPMLog {
private:
	persistent_ptr<LSNmap> lsnPptr;
    persistent_ptr<GSNmap> gsnPptr;	
	p<uint64_t> highest_gsn;
public:
	p<pool<root>> pop;
    cppPMLog(pool<root>);
    ~cppPMLog();
    void restartMaps();
    int Append(const char* record, uint64_t lsn);
    int Commit(uint64_t lsn, uint64_t gsn) ;
    uint64_t Read(uint64_t gsn, char* storage);
    int Trim(uint64_t gsn);
};

void* cppStartUp();
void cppFinalize(persistent_ptr<cppPMLog> cppLog);