CC=g++
OPT=-static -O3
#OPT=-static -g
LIBS=-I../../lib/cnpy -L../../lib/cnpy -lcnpy

all: npy-consolidate band-to-matrix-npy npy-vec-to-hiq-1hot

.PHONY: clean
clean:
	rm npy-consolidate band-to-matrix-npy npy-vec-to-hiq-1hot

npy-consolidate: npy-consolidate.cpp
	$(CC) $(OPT) npy-consolidate.cpp -o npy-consolidate $(LIBS)

band-to-matrix-npy: band-to-matrix-npy.cpp
	$(CC) $(OPT) band-to-matrix-npy.cpp -o band-to-matrix-npy $(LIBS)

npy-vec-to-hiq-1hot: npy-vec-to-hiq-1hot.cpp
	$(CC) $(OPT) npy-vec-to-hiq-1hot.cpp -o npy-vec-to-hiq-1hot $(LIBS)



