
CC=g++
CC_FLAGS=-O3 -std=c++11 -msse4.2 -lhts
CC_LIBS=-L$(HOME)/lib -L$(HOME)/lib/htslib
CC_INC=-I$(HOME)/lib/htslib/htslib

all: tile-assembly

.PHONY: clean
clean:
	rm -f *.o
	rm -f tile-assembly

tile-assembly: *.cpp
	$(CC) $(CC_INC) $^ -o tile-assembly $(CC_LIBS) $(CC_FLAGS) $(CC_LIBS)/libhts.a  


