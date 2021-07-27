CXX:=g++
#OPT=-O3
CXXFLAGS+=-g
LIBS=-L../../lib/asmukk -I../../lib/asmukk ../../lib/asmukk/asm_ukk.a

all: tile-liftover tile-assembly-index

.PHONY: clean
clean:
	rm -f tile-liftover tile-assembly-index

tile-liftover: tile-liftover.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) tile-liftover.cpp -o tile-liftover $(LIBS)

tile-assembly-index: tile-assembly-index.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) tile-assembly-index.cpp -o tile-assembly-index
