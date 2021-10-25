main: main.o libstorage.so
	g++ main.o -o main -L. -lstorage -lpmemobj -ltbb

main.o: main.cpp
	g++ -c main.cpp -L. -lstorage

libstorage.so: PMLog.o cPMLog.o LogCache.o
	g++ PMLog.o cPMLog.o LogCache.o -fPIC -shared -o libstorage.so

PMLog.o: PMLog.cpp PMLog.hpp
	g++ PMLog.cpp -c -fPIC -L. -lpmemobj -pthread -o PMLog.o
	
cPMLog.o: cPMLog.cpp PMLog.h
	g++ cPMLog.cpp -c -fPIC -o cPMLog.o
	
LogCache.o: LogCache.cpp LogCache.hpp
	g++ LogCache.cpp -c -fPIC -o LogCache.o
	
clean: 
	rm libstorage.so PMLog.o cPMLog.o LogCache.o
