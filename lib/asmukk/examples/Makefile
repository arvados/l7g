CC=gcc

all: dp mktest dp_dna mkseq

dp: dp.c
	$(CC) $< -o $@

dp_dna: dp_dna.c
	$(CC) $< -o $@

mktest: mktest.c
	$(CC) $< -o $@

mkseq: mkseq.c
	$(CC) $< -o $@

test: dp dp_dna mktest mkseq
	./check_ukk.sh
	./check_ukk_dna.sh

clean:
	rm -f dp mktest dp_dna mkseq *~
