CC=g++
OPT=-static -O3
LIBS=-I../../lib/cnpy -L../../lib/cnpy -lcnpy

npy-consolidate: npy-consolidate.cpp
	$(CC) $(OPT) npy-consolidate.cpp -o npy-consolidate $(LIBS)

band-to-matrix-npy: band-to-matrix-npy.cpp
	$(CC) $(OPT) band-to-matrix-npy.cpp -o band-to-matrix-npy $(LIBS)

npy-vec-to-hiq-1hot: npy-vec-to-hiq-1hot.cpp
	$(CC) $(OPT) npy-vec-to-hiq-1hot.cpp -o npy-vec-to-hiq-1hot $(LIBS)
