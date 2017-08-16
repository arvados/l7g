CC=gcc
CPP=g++
CFLAGS=-I. -L. -O3

all: asm_ukk.o asm_ukk asm_ukk.a asm_ukk_dna

#asm_ukk.o: asm_ukk.c asm_ukk3.c asm_ukk.h
asm_ukk3.o: asm_ukk3.c asm_ukk.h
	$(CC) -c -o $@ asm_ukk3.c $(CFLAGS)

asm_ukk.o: asm_ukk.c asm_ukk.h
	$(CC) -c -o $@ asm_ukk.c $(CFLAGS)

asm_ukk.a: asm_ukk.o asm_ukk3.o
	ar rcs $@ asm_ukk.o asm_ukk3.o

asm_ukk: asm_ukk.a asm_ukk_main.cpp
	$(CPP) -o $@ asm_ukk_main.cpp $(CFLAGS) asm_ukk.a

asm_ukk_dna: asm_ukk.a asm_ukk_dna.cpp
	$(CPP) -o $@ asm_ukk_dna.cpp $(CFLAGS) asm_ukk.a

clean:
	rm -f *.o *.a *~ asm_ukk asm_ukk_dna
