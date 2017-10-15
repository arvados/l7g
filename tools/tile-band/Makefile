CC=g++
CC_FLAGS=-O3 -std=c++11 -msse4.2
#CC_LIB=
#CC_INC=

src=$(wildcard *.cpp)
obj=$(src:.cpp=.o)

tile-band: $(src)
	$(CC) $(CC_FLAGS) $^ -o tile-band

.PHONY: clean
clean:
	rm -f *.o
	rm -f tile-band

