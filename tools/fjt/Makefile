#OPT=-g
OPT=-O3
LIBS=-lcrypto -lz
fjt: fjt.cpp cJSON.c sglf.cpp sglf.hpp twoBit.hpp twoBit.cpp
	g++ ${OPT} fjt.cpp cJSON.c sglf.cpp twoBit.cpp -o fjt ${LIBS}
