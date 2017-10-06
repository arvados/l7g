FROM arvados/jobs:latest


USER root

RUN apt-get update && \
  apt-get install -y build-essential golang libssl-dev \
    less wget \
    libncurses5-dev autoconf libbz2-dev liblzma-dev \
    cmake tabix libhts-dev

RUN mkdir -p /root/git && \
  cd /root/git && \
  git clone https://github.com/curoverse/l7g && \
  git clone https://github.com/simongog/sdsl-lite && \
  cd /root/git/sdsl-lite && ./install.sh /usr/local && \

  cd /root/git && \
  git clone https://github.com/samtools/htslib && \
  cd htslib && \
  autoheader && autoconf && \
  ./configure && make && make install && \

  cd /root/git && \
  cd /root/git && \
  git clone https://github.com/samtools/samtools && \
  cd /root/git/samtools && \
  autoheader && autoconf -Wno-syntax && \
  ./configure && make && make install && \

  cd /root/git/l7g/tools/fjt && make && \
  cp /root/git/l7g/tools/fjt/fjt /usr/local/bin && \

  cd /root/git/l7g/tools/cgft && make && \
  cp /root/git/l7g/tools/cgft/cgft /usr/local/bin && \

  cd /root/git/l7g/tools/tile-assembly && make && \
  cp /root/git/l7g/tools/tile-assembly/tile-assembly /usr/local/bin && \

  cd /root/git/l7g/tools/tile-lib && go build merge-tilelib.go && \
  cp /root/git/l7g/tools/tile-lib/merge-tilelib /usr/local/bin  && \

  cd /root/git/l7g/tools/tagset/src && make && \
  cp /root/git/l7g/tools/tagset/src/tagsetFa /usr/local/bin && \
  cp /root/git/l7g/tools/tagset/choose_tagset_startpos0_vestigial.py /usr/local/bin && \
  cp /root/git/l7g/tools/tagset/choose_tagset_startpos0.py /usr/local/bin && \

  cd /root/git/l7g/tools/l7g && go build l7g.go && \
  cp /root/git/l7g/tools/l7g/l7g /usr/local/bin && \

  cd /root/git/l7g/lib/cnpy && cmake . && make && \
  cd /root/git/l7g/tools/tile-npy && make && \
  cp /root/git/l7g/tools/tile-npy/band-to-matrix-npy /usr/local/bin  && \
  cp /root/git/l7g/tools/tile-npy/npy-consolidate /usr/local/bin  && \
  cp /root/git/l7g/tools/tile-npy/npy-vec-to-hiq-1hot /usr/local/bin  && \
  cp /root/git/l7g/tools/tile-npy/create-list /usr/local/bin  && \

  cd /root/git/l7g/tools/glft && make && \
  cp /root/git/l7g/tools/glft/glft /usr/local/bin && \

  cd /root/git/l7g/tools/cglf-tools && \
  GOPATH=/root/.go go get github.com/codegangsta/cli && \
  GOPATH=/root/.go ./cmp_fastj2cglf.sh && \
  cp /root/git/l7g/tools/cglf-tools/fastj2cgflib /usr/local/lib && \

  cp /root/git/l7g/tools/misc/refstream /usr/local/lib && \

  cd /root/git/l7g/go/pasta/pasta && \
  GOPATH=/root/.go go get github.com/curoverse/l7g/go/autoio && \
  GOPATH=/root/.go go get github.com/curoverse/l7g/go/memz && \
  GOPATH=/root/.go go get github.com/curoverse/l7g/go/pasta && \
  GOPATH=/root/.go go get github.com/curoverse/l7g/go/pasta/gvcf && \
  GOPATH=/root/.go go get github.com/curoverse/l7g/go/sloppyjson && \
  GOPATH=/root/.go go build && \
  cp /root/git/l7g/go/pasta/pasta/pasta /usr/local/bin && \

  cd /root/git/l7g/go/fjtools && \
  GOPATH=/root/.go go build fjfilter.go && \
  cp /root/git/l7g/go/fjtools/fjfilter /usr/local/bin

RUN apt-get install -y bc

RUN cd /usr/local/bin && \
  wget http://hgdownload.cse.ucsc.edu/admin/exe/linux.x86_64/bigWigToBedGraph && \
  chmod a+x bigWigToBedGraph

USER crunch
