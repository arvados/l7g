CXX:=g++
CXXFLAGS+=-O3
#OPT=-static -g
LIBS=-I../../lib/cnpy -L../../lib/cnpy -lcnpy

all: npy-consolidate band-to-matrix-npy npy-vec-to-hiq-1hot

.PHONY: clean
clean:
	rm -f npy-consolidate band-to-matrix-npy npy-vec-to-hiq-1hot || true

npy-consolidate: npy-consolidate.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) npy-consolidate.cpp -o npy-consolidate $(LIBS)

band-to-matrix-npy: band-to-matrix-npy.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) band-to-matrix-npy.cpp -o band-to-matrix-npy $(LIBS)

npy-vec-to-hiq-1hot: npy-vec-to-hiq-1hot.cpp
	$(CXX) $(LDFLAGS) $(CPPFLAGS) $(CXXFLAGS) npy-vec-to-hiq-1hot.cpp -o npy-vec-to-hiq-1hot $(LIBS)



