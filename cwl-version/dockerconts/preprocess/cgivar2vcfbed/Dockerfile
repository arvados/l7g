FROM arvados/jobs
MAINTAINER Jiayong Li <jli@curii.com>

USER root

RUN apt-get update -q

RUN apt-get install -qy build-essential wget cmake zlib1g-dev git

# Installing cgatools 1.8.0

RUN wget https://sourceforge.net/projects/cgatools/files/1.8.0/cgatools-1.8.0.1-linux_binary-x86_64.tar.gz && \
    tar -xzvf cgatools-1.8.0.1-linux_binary-x86_64.tar.gz && \
    cp cgatools-1.8.0.1-linux_binary-x86_64/bin/cgatools /usr/local/bin    

WORKDIR /

# Getting HTSLIB for tabix/bgzip

RUN wget https://github.com/samtools/htslib/releases/download/1.3.1/htslib-1.3.1.tar.bz2 && \
    tar -xjvf htslib-1.3.1.tar.bz2 && \
    cd htslib-1.3.1 && ./configure && make && make install

WORKDIR /

# Installing bedtools v2.25.0 and dependencies

RUN wget https://github.com/arq5x/bedtools2/releases/download/v2.25.0/bedtools-2.25.0.tar.gz && \
    tar -xzvf bedtools-2.25.0.tar.gz && \
    cd bedtools2/ && make && cp /bedtools2/bin/bedtools /usr/local/bin

WORKDIR /

# Installing gvcf_regions

RUN git clone https://github.com/lijiayong/gvcf_regions
