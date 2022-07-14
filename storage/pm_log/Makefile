CCX=g++
CXX_FLAGS= -Wall -g -O3
WORKING_DIR=/home/${USER}/PMLog-bench
INCLUDE_PATH=-I${WORKING_DIR}/pmdk/local/usr/include -I${WORKING_DIR}/libpmemobj-cpp/include
LD_FLAGS=-L${WORKING_DIR}/pmdk/local/lib -L.
LIBS=-lfmt -lstorage -lpmemobj -ltbb -lpmem -pthread -lgflags

.PHONY = all clean 
all: clean pmem_bench main

pmem_bench: libstorage.so
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) $@.cpp -o $@ ${LD_FLAGS} ${LIBS}

main: main.o libstorage.so
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) $@.o -o $@ ${LD_FLAGS} ${LIBS}

main.o: main.cpp
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) -c main.cpp  ${LD_FLAGS} ${LIBS}

libstorage.so: PMLog.o cPMLog.o LogCache.o
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) PMLog.o cPMLog.o LogCache.o -fPIC -shared -o libstorage.so

PMLog.o: PMLog.cpp PMLog.hpp
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) PMLog.cpp -c -fPIC ${LD_FLAGS} ${LIBS} -o PMLog.o
	
cPMLog.o: cPMLog.cpp PMLog.h
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) cPMLog.cpp -c -fPIC -o cPMLog.o
	
LogCache.o: LogCache.cpp LogCache.hpp
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) LogCache.cpp -c -fPIC -o LogCache.o
	
clean: 
	-rm -f libstorage.so PMLog.o cPMLog.o LogCache.o main.o main
