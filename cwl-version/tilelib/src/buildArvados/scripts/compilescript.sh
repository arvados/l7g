
set -e 
mkdir -p $HOME/go
export GOPATH=$HOME/go

#go get -u 'github.com/curoverse/l7g/tools/cglf-tools/fast2cgflib'
#cp $HOME/go/bin/fasat2cgflib ../dest

go get -u 'github.com/curoverse/l7g/tools/tilelib/merge-tilelib'
cp $HOME/go/bin/merge-tilelib ../dest

