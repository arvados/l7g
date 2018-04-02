set -e 
mkdir -p $HOME/go
export GOPATH=$HOME/go


wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip


(cd '../src/l7g/l7g-master/tools/fjt'; make;) # cp -p cgft ${DEST}/fjt)
cp -p ../src/l7g/l7g-master/tools/fjt/fjt ../dest/fjt


(cd '../src/l7g/l7g-master/tools/fjcsv2sglf'; make;) # cp -p cgft ${DEST}/fjcsv2sglf)
cp -p ../src/l7g/l7g-master/tools/fjcsv2sglf/fjcsv2sglf ../dest/fjcsv2sglf


(cd '../src/l7g/l7g-master/tools/tile-lib'; make;) # cp -p cgft ${DEST}/merge-sglf)
cp -p ../src/l7g/l7g-master/tools/tile-lib/merge-sglf ../dest/merge-sglf

