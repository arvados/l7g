CC=g++
#OPT=-g -lcrypto
OPT=-O3
LIB=-lcrypto

fjcsv2sglf: twoBit.cpp fjcsv2sglf.cpp
	${CC} ${OPT} fjcsv2sglf.cpp twoBit.cpp -o fjcsv2sglf ${LIB}

twobit_test: twoBit.cpp
	${CC} ${OPT} -DTWOBIT_TEST twoBit.cpp -o twobit_test ${LIB}


