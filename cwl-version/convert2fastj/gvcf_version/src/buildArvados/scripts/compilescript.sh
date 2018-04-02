set -e 
mkdir -p $HOME/go 
export GOPATH=$HOME/go

go get -u 'github.com/curoverse/l7g/tools/l7g'
go get -u 'github.com/curoverse/l7g/go/pasta/pasta'

wget 'https://raw.githubusercontent.com/curoverse/l7g/master/tools/misc/refstream' 

cp $HOME/go/bin/l7g ../dest
cp $HOME/go/bin/pasta ../dest

chmod +x refstream
mv refstream ../dest

wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip

ASM_UKK=../src/l7g/l7g-master/lib/asmukk

g++ -g ../src/l7g/l7g-master/tools/which-ref/which-ref.cpp ${ASM_UKK}/asm_ukk.c -o ../dest/which-ref -I${ASM_UKK}

CC_FLAGS='-O3 -std=c++11 -msse4.2 -lhts -lpthread -lz -lm -llzma -lbz2'

g++  ../src/l7g/l7g-master/tools/tile-assembly/tile-assembly.cpp -o ../dest/tile-assembly $CC_FLAGS
