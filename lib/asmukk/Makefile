CC:=gcc
CXX:=g++
CFLAGS+=-I. -L. -O3
CXXFLAGS+=-I. -L. -O3
#CFLAGS=-I. -L. -g

all: asm_ukk.o asm_ukk asm_ukk.a asm_ukk_dna avm_ukk3.o aim_ukk

.PHONY: clean
clean:
	rm -f *.o *.a *~ asm_ukk asm_ukk_dna aim_ukk

#asm_ukk.o: asm_ukk.c asm_ukk3.c asm_ukk.h
asm_ukk3.o: asm_ukk3.c asm_ukk.h
	$(CC) $(CPPFLAGS) $(CFLAGS) -c -o $@ asm_ukk3.c

asm_ukk.o: asm_ukk.c asm_ukk.h
	$(CC) $(CPPFLAGS) $(CFLAGS) -c -o $@ asm_ukk.c

asm_ukk.a: asm_ukk.o asm_ukk3.o avm_ukk3.o
	ar rcs $@ asm_ukk.o asm_ukk3.o avm_ukk3.o

asm_ukk: asm_ukk.a asm_ukk_main.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) -o $@ asm_ukk_main.cpp asm_ukk.a

asm_ukk_dna: asm_ukk.a asm_ukk_dna.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) -o $@ asm_ukk_dna.cpp asm_ukk.a

avm_ukk3.o: avm_ukk3.c
	$(CC) $(CPPFLAGS) $(CFLAGS) -c -o $@ avm_ukk3.c

aim_ukk: avm_ukk3.o aim_ukk.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) -o $@ aim_ukk.cpp asm_ukk.a

