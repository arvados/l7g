#OPT=-g
CXXFLAGS+=-O3
LIBS=-lcrypto -lz
fjt: fjt.cpp cJSON.c sglf.cpp sglf.hpp twoBit.hpp twoBit.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) fjt.cpp cJSON.c sglf.cpp twoBit.cpp -o fjt ${LIBS}
