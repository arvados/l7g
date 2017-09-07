CC=g++
#CC_FLAGS=-g -std=c++11 -msse4.2
CC_FLAGS=-O3 -std=c++11 -msse4.2
CC_LIB=-lsdsl -L$(HOME)/lib
CC_INC=-I$(HOME)/include

src=$(wildcard *.cpp)
obj=$(src:.cpp=.o)


%.o: %.cpp
	$(CC) -c -o $@ $^ $(CC_FLAGS) $(CC_LIB) $(CC_INC)

cgft: $(obj)
	$(CC) $(CC_FLAGS) $^ -o cgft $(CC_LIB) $(CC_INC)
#	$(CC) $(CC_FLAGS) band-to-enc.o band-to-enc-io.o cgft_container.o cgft_tilemap.o cgft_io.o cgft_alloc.o cgft_encode.o main.cpp -o cgft $(CC_LIB) $(CC_INC)

.PHONY: clean
clean:
	rm -f *.o
	rm -f cgft

