libstorage.so: PMLog.o cPMLog.o
	g++ PMLog.o cPMLog.o -fPIC -shared -o libstorage.so

PMLog.o: PMLog.cpp PMLog.hpp
	g++ PMLog.cpp -c -fPIC -L. -lpmemobj -o PMLog.o
	
cPMLog.o: cPMLog.cpp PMLog.h
	g++ cPMLog.cpp -c -fPIC -o cPMLog.o
	
clean: 
	rm libstorage.so PMLog.o cPMLog.o
