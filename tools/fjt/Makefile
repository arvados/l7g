#OPT='-g'
OPT='-O3'
LIBS='-lcrypto'
fjt: fjt.cpp cJSON.c sglf.cpp sglf.hpp
	g++ ${OPT} fjt.cpp cJSON.c sglf.cpp -o fjt ${LIBS}
