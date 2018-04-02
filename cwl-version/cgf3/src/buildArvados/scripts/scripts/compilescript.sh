set -e 
mkdir -p $HOME/go
export GOPATH=$HOME/go

wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip

LIB='-lcrypto'
FJTDIR='../src/l7g/l7g-master/tools/fjt'

#echo $DEST

(cd '../src/l7g/l7g-master/tools/fjt'; make;) # cp -p cgft ${DEST}/fjt)
cp -p ../src/l7g/l7g-master/tools/fjt/fjt ../dest/fjt

(cd '../src/l7g/l7g-master/tools/cgft'; make;) # cp -p cgft ${DEST}/cgft)
cp -p ../src/l7g/l7g-master/tools/cgft/cgft ../dest/cgft
