#include "PMLog.h"
#include <iostream>
#include <stdlib.h>
#include <fmt/printf.h>
#include <gflags/gflags.h>
#include <thread>
#include <vector>

using GFLAGS_NAMESPACE::ParseCommandLineFlags;
using GFLAGS_NAMESPACE::RegisterFlagValidator;
using GFLAGS_NAMESPACE::SetUsageMessage;

DEFINE_int64(reads, 500, "Percentage of read operations to do.");
DEFINE_int32(threads, 1, "Number of concurrent threads to run.");

DEFINE_int32(record_size, 100, "Size of each record to be appended");
DEFINE_int32(key_size, 16, "Size of each key");
DEFINE_int32(nb_ops, 1e6, "Number of total operations to be executed.");

enum DistributionType : unsigned char {
	kFixed = 0,
	kUniform,
	kNormal    
};

#if 0
static enum DistributionType FLAGS_distribution_type = kFixed;

static enum DistributionType StringToDistributionType(const char* ctype) {
	if (!strcasecmp(ctype, "fixed"))                                        
		return kFixed;
	else if (!strcasecmp(ctype, "uniform"))
		return kUniform;
	else if (!strcasecmp(ctype, "normal"))
		return kNormal;

	fmt::print("[{}] Cannot parse distribution type {}\n", __func__, ctype);
	return kFixed;  // default value
}
#endif

class KeyGenerator {
	public:
		KeyGenerator(uint64_t num) : num_(num) {}

		uint64_t Next() {
			return rand() % num_;
		}

	private:
		const uint64_t num_;
};


std::unique_ptr<char[]> random_str(size_t sz) {
	return std::make_unique<char[]>(sz);
}



static void thread_func(int&& thread_id) {
	auto record = random_str(FLAGS_record_size);
	PMLog log = startUp_idx(thread_id);
	KeyGenerator gen(FLAGS_nb_ops/2);

	std::unique_ptr<uint64_t> next_gsn = std::make_unique<uint64_t>();

	// loading phase
	for (int i = 0; i < FLAGS_nb_ops/2; i++) {

		cAppend(log, record.get(), i);
		cCommit(log, i, i);
	}

	for (int i = 0; i < FLAGS_nb_ops; i++) {
		if (rand()%1000 > FLAGS_reads) {
			cAppend(log, record.get(), i);
			cCommit(log, i, i);
		}
		else
			cRead(log, i, next_gsn.get());
		// std::cout << cRead(log, i, next_gsn) << std::endl;
		if (i%5125 == 0)
			fmt::print("[{}] {}\n", __func__, i);
	}

	fmt::print("[{}] thread={} finished ...\n", __func__, thread_id);
	finalize(log);

	fmt::print("[{}] thread={} finished ...\n", __func__, thread_id);
}

int main(int args, char* argv[]) {
	ParseCommandLineFlags(&args, &argv, true);
	fmt::print("[{}] reads={}\tthreads={}\trecord_size={}\tkey_size={}\tnb_ops={}\n", __func__, FLAGS_reads, FLAGS_threads, FLAGS_record_size, FLAGS_key_size, FLAGS_nb_ops);
	std::vector<std::thread> threads;
	for (auto i = 0ULL; i < FLAGS_threads; i++)
		threads.emplace_back(thread_func, i);

	for (auto& thread: threads)
		thread.join();

	fmt::print("[{}] finished\n", __func__);

}
