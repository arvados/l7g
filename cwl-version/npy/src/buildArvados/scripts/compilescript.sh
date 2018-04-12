# compiling: band-to-matrix,npy-vect-to-hiq-1hot,npy-consolidate   
export liblocation=$1
mkdir -p ~/go
export GOPATH=~/go

wget https://github.com/curoverse/cgf/archive/master.zip
mkdir -p ../src/cgf
unzip -d ../src/cgf master.zip
rm master.zip
 
wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip

g++ -static -o ../dest/band-to-matrix-npy ../src/l7g/l7g-master/tools/tile-npy/band-to-matrix-npy.cpp -L$liblocation -lcnpy -I$liblocation 

g++ -static -o ../dest/npy-vec-to-hiq-1hot ../src/l7g/l7g-master/tools/tile-npy/npy-vec-to-hiq-1hot.cpp -L$liblocation -lcnpy -I$liblocation

g++ -static -o ../dest/npy-consolidate ../src/l7g/l7g-master/tools/tile-npy/npy-consolidate.cpp -L$liblocation -lcnpy -I$liblocation

(cd '../src/l7g/l7g-master/tools/cgft'; make;) # cp -p cgft ${DEST}/cgft)
cp -p ../src/l7g/l7g-master/tools/cgft/cgft ../dest/cgft

