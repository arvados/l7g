#!/bin/bash

g++ -O3 -Llib/cnpy -lcnpy -Ilib/cnpy process-npy.cpp -o process-npy -Llib/cnpy -lcnpy -Ilib/cnpy

exit

g++ -O3 -Llib/cnpy -lcnpy -Ilib/cnpy process-npy.cpp -o process-npy -Llib/cnpy -lcnpy -Ilib/cnpy

exit

g++ -g -Llib/cnpy -lcnpy -Ilib/cnpy band-to-matrix-npy.cpp -o band-to-matrix-npy -Llib/cnpy -lcnpy -Ilib/cnpy
g++ -g -Llib/cnpy -lcnpy -Ilib/cnpy npy-to-hiq.cpp -o npy-to-hiq -Llib/cnpy -lcnpy -Ilib/cnpy
