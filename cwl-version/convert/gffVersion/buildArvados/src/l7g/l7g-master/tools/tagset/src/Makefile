HTSDIR=${HOME}/lib/htslib
HTSLIB=${HOME}/lib/htslib
OPT=-O3

tiletagFa: tagset-fa.cpp
	g++ ${OPT} -I ${HTSDIR}/htslib -lhts tagset-fa.cpp -o tagsetFa -L ${HTSLIB} -lhts
