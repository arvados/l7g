CPP=g++
#OPT=-O3
OPT=-g
LIBS=-L../../lib/asmukk -I../../lib/asmukk ../../lib/asmukk/asm_ukk.a

all: tile-liftover

.PHONY: clean
clean:
	rm -f tile-liftover

tile-liftover: tile-liftover.cpp
	g++ $(OPT) tile-liftover.cpp -o tile-liftover $(LIBS)
