wget https://github.com/curoverse/l7g/archive/master.zip
mkdir -p ../src/l7g
unzip -o -d ../src/l7g master.zip
rm master.zip

g++ -O3 -o ../dest/cleanvcf ../src/l7g/l7g-master/tools/misc/cleanvcf.cpp
 
