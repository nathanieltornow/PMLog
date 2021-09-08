#include <libpmemobj++/persistent_ptr.hpp>
#include <libpmemobj++/container/concurrent_hash_map.hpp>
#include <cstdint>
#include <string>

#define KEY_SIZE 16
using namespace pmem::obj;

class PersistentString {
private:
	char array[KEY_SIZE];
public:
	std::string data() const;	
	PersistentString(std::string* s);
};

using PString = struct PersistentString;
using LSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;
using GSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;
using ColorMap = concurrent_hash_map<p<uint32_t>, persistent_ptr<GSNmap>>;

class cppPMLog {
private:
	persistent_ptr<LSNmap> lsnPptr;
    persistent_ptr<ColorMap> colorPptr;
	pool_base pop;
public:
    cppPMLog(pool_base);
    ~cppPMLog();
    void restartMaps();
    int Append(std::string record, uint64_t lsn);
    int Commit(uint64_t lsn, uint32_t color, uint64_t gsn) ;
    std::string Read(uint32_t color, uint64_t gsn);
    int Trim(uint32_t color, uint64_t gsn);
};

void* cppStartUp();