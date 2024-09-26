CXXFLAGS+=-O3
LDLIBS=-lcrypto
tileband-hash: sglf.cpp main.cpp tileband-hash.cpp twoBit.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) $^ -o $@ $(LDLIBS)

