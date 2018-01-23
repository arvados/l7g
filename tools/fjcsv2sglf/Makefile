CC=g++
#OPT=-g -lcrypto
OPT=-O3
OPT_DEBUG=-g
LIB=-lcrypto

fjcsv2sglf: twoBit.cpp fjcsv2sglf.cpp
	${CC} ${OPT} fjcsv2sglf.cpp twoBit.cpp -o fjcsv2sglf ${LIB}

fjcsv2sglf_debug: twoBit.cpp fjcsv2sglf.cpp
	${CC} ${OPT_DEBUG} -DDEBUG_RUN fjcsv2sglf.cpp twoBit.cpp -o fjcsv2sglf_debug ${LIB}

twobit_test: twoBit.cpp
	${CC} ${OPT} -DTWOBIT_TEST twoBit.cpp -o twobit_test ${LIB}


