
CC=g++
CC_FLAGS=-O3 -std=c++11 -msse4.2 -lhts

all: tile-assembly

.PHONY: clean
clean:
	rm -f *.o
	rm -f tile-assembly

tile-assembly: *.cpp
	$(CC)  $^ -o tile-assembly $(CC_FLAGS)


