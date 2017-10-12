set -e 
mkdir -p $HOME/go 
export GOPATH=$HOME/go

wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip

#LIB='-lcrypto'
#FJTDIR='../src/l7g/l7g-master/tools/fjt'

#echo $DEST

#g++ -g ${FJTDIR}/fjt.cpp ${FJTDIR}/cJSON.c ${FJTDIR}/sglf.cpp -o ../dest/fjt ${LIB} 

#(cd '../src/l7g/l7g-master/tools/cgft'; make;) # cp -p cgft ${DEST}/cgft)
#cp -p ${SRC}/l7g/l7g-master/tools/cgft/cgft ${DEST}/cgft

