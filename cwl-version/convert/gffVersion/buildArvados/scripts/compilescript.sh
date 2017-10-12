set -e 
mkdir -p ~/go
export GOPATH=~/go
go get 'github.com/curoverse/l7g/tools/l7g'
go get 'github.com/curoverse/l7g/go/pasta/pasta'
go get 'github.com/boopathi/numberify'

wget 'https://raw.githubusercontent.com/curoverse/l7g/master/tools/misc/refstream' 
cp ~/go/bin/l7g ../dest
cp ~/go/bin/pasta ../dest
cp ~/go/bin/numberify ../dest

chmod +x refstream
mv refstream ../dest

wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip

ASM_UKK=../src/l7g/l7g-master/lib/asmukk

g++ -g ../src/l7g/l7g-master/tools/which-ref/which-ref.cpp ${ASM_UKK}/asm_ukk.c -o ../dest/which-ref -I${ASM_UKK}
