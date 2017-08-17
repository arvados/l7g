# compiling: band-to-matrix-npy.cpp  npy-to-hiq.cpp  process-npy.cpp  
export liblocation=$1

cd .

#g++ ./src/band-to-matrix-npy.cpp -L/data-sdd/cwl_tiling/npy/lib/cnpy -lcnpy -I/data-sdd/cwl_tiling/npy/lib/cnpy -o ./dest/band-to-matrix
#g++ ./src/npy-to-hiq.cpp -L/data-sdd/cwl_tiling/npy/lib/cnpy -lcnpy -I/data-sdd/cwl_tiling/npy/lib/cnpy -o ./dest/npy-to-hiq
#g++ ./src/process-npy.cpp -L/data-sdd/cwl_tiling/npy/lib/cnpy -lcnpy -I/data-sdd/cwl_tiling/npy/lib/cnpy -o ./dest/process-npy

g++ -o ./dest/band-to-matrix-npy ./src/band-to-matrix-npy.cpp -L$liblocation -lcnpy -I$liblocation 

g++ -o ./dest/npy-to-hiq ./src/npy-to-hiq.cpp -L$liblocation -lcnpy -I$liblocation 

g++ -o ./dest/process-npy ./src/process-npy.cpp -L$liblocation -lcnpy -I$liblocation

g++ -o ./dest/npy-vec-to-hiq-1hot ./src/npy-vec-to-hiq-1hot.cpp -L$liblocation -lcnpy -I$liblocation

