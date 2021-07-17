
CXX:=g++
CXXFLAGS+=-O3 -std=c++11
LDLIBS+=-lhts

all: tile-assembly

.PHONY: clean
clean:
	rm -f *.o
	rm -f tile-assembly

tile-assembly: *.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS)  $^ -o tile-assembly $(LDLIBS)


