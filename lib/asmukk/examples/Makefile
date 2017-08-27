CC=gcc
GPP=g++

all: dp mktest dp_dna mkseq mkiarray

dp: dp.c
	$(CC) $< -o $@

dp_dna: dp_dna.c
	$(GPP) $< -o $@

dpi: dpi.c
	$(GPP) $< -o $@

mktest: mktest.c
	$(CC) $< -o $@

mkseq: mkseq.c
	$(CC) $< -o $@

mkiarray: mkiarray.c
	$(GPP) $< -o $@

test: dp dp_dna mktest mkseq
	./check_ukk.sh
	./check_ukk_dna.sh

.PHONY: clean
clean:
	rm -f dp mktest dp_dna mkseq *~
