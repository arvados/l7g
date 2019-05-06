#
# Copyright (C) The Arvados Authors. All rights reserved.
#
# SPDX-License-Identifier: AGPL-3.0
#

FROM arvados/jobs:latest

USER root

RUN apt-get update && \
  apt-get install -y build-essential golang libssl-dev git bedtools \
    less wget libncursesw5-dev vcftools \
    libncurses5-dev libcurl4-gnutls-dev autoconf libbz2-dev liblzma-dev \
    cmake tabix zlib1g-dev libhts-dev parallel python-numpy


RUN mkdir -p /root/git && \
  cd /root/git && \
  git clone https://github.com/curoverse/l7g && \
  git clone https://github.com/simongog/sdsl-lite && \
  cd /root/git/sdsl-lite && ./install.sh /usr/local && \

  cd /usr/bin && \
  wget https://github.com/samtools/htslib/releases/download/1.9/htslib-1.9.tar.bz2 && \
  tar -vxjf htslib-1.9.tar.bz2 && \
  cd htslib-1.9 && \
  ./configure && \
  make && \
  make install && \

  cd /usr/bin && \
  wget https://github.com/samtools/samtools/releases/download/1.9/samtools-1.9.tar.bz2 && \
  tar -vxjf samtools-1.9.tar.bz2 && \
  cd samtools-1.9 && \
  ./configure && \
  make && \
  make install && \

  cd /root/git/l7g/tools/misc && make && \
  cp /root/git/l7g/tools/misc/cleanvcf /usr/local/bin && \

  cd /root/git/l7g/tools/fjt && make && \
  cp /root/git/l7g/tools/fjt/fjt /usr/local/bin && \

  cd /root/git/l7g/tools/fjcsv2sglf && make && \
  cp /root/git/l7g/tools/fjcsv2sglf/fjcsv2sglf /usr/local/bin && \

  cd /root/git/l7g/tools/which-ref && make && \
  cp /root/git/l7g/tools/which-ref/which-ref /usr/local/bin && \

  cd /root/git/l7g/tools/cgft && make && \
  cp /root/git/l7g/tools/cgft/cgft /usr/local/bin && \

  cd /root/git/l7g/tools/tile-assembly && make && \
  cp /root/git/l7g/tools/tile-assembly/tile-assembly /usr/local/bin && \

  cd /root/git/l7g/tools/tile-lib && go build merge-tilelib.go && \
  make all && \
  cp /root/git/l7g/tools/tile-lib/merge-tilelib /usr/local/bin  && \
  cp /root/git/l7g/tools/tile-lib/merge-sglf /usr/local/bin  && \
  cp /root/git/l7g/tools/tile-lib/sglf-check /usr/local/bin  && \

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

  cd /root/git/l7g/tools/tileband-hash && make && \
  cp /root/git/l7g/tools/tileband-hash/tileband-hash /usr/local/bin && \

  cd /root/git/l7g/tools/vcfbed2homref && make && \
  cp /root/git/l7g/tools/vcfbed2homref/vcfbed2homref /usr/local/bin && \

  cd /root/git/l7g/tools/cglf-tools && \
  GOPATH=/root/.go go get github.com/codegangsta/cli && \
  GOPATH=/root/.go ./cmp_fastj2cglf.sh && \
  cp /root/git/l7g/tools/cglf-tools/fastj2cgflib /usr/local/bin && \

  cp /root/git/l7g/tools/misc/refstream /usr/local/bin && \

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

RUN apt-get install -y bc file

RUN  cd /root/git/l7g/lib/asmukk && make && \
  cd /root/git/l7g/tools/tile-liftover && make && \
  cp /root/git/l7g/tools/tile-liftover/tile-liftover /usr/local/bin && \
  cp /root/git/l7g/tools/tile-liftover/tile-assembly-index /usr/local/bin

RUN cd /usr/local/bin && \
  wget http://hgdownload.cse.ucsc.edu/admin/exe/linux.x86_64/bigWigToBedGraph && \
  chmod a+x bigWigToBedGraph
