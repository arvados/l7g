CXX:=g++
CXXFLAGS+=-O3
OPT_DEBUG=-g

merge-sglf: merge-sglf.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) merge-sglf.cpp -o merge-sglf

sglf-check: sglf-check.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) sglf-check.cpp -o sglf-check -lcrypto

merge-sglf_debug: merge-sglf.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) ${OPT_DEBUG} merge-sglf.cpp -o merge-sglf_debug

all: merge-sglf sglf-check merge-sglf_debug
