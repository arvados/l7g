CC=g++
OPT=-O3
OPT_DEBUG=-g

merge-sglf: merge-sglf.cpp
	${CC} ${OPT} merge-sglf.cpp -o merge-sglf

sglf-check: sglf-check.cpp
	${CC} ${OPT} sglf-check.cpp -o sglf-check -lcrypto

merge-sglf_debug: merge-sglf.cpp
	${CC} ${OPT_DEBUG} merge-sglf.cpp -o merge-sglf_debug
