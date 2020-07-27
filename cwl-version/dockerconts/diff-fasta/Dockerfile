FROM debian:buster-slim
MAINTAINER Jiayong Li <jli@curii.com>
USER root

# Install necessary dependencies

RUN apt-get update --fix-missing -qy

RUN apt-get install -qy wget \
        build-essential \
        cmake \
        zlib1g-dev \
        libbz2-dev \
        liblzma-dev \
        libncurses5-dev \
        libncursesw5-dev \
        gcc \
        make \
        python \
        git

# Installing go 1.14

RUN wget https://golang.org/dl/go1.14.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.14.4.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

# Getting lightning (diff-fasta)

RUN go get -v github.com/arvados/lightning

ENV PATH $PATH:/root/go/bin

# Getting HTSLIB 1.9 for tabix/bgzip

RUN wget https://github.com/samtools/htslib/releases/download/1.9/htslib-1.9.tar.bz2 && tar -xjvf htslib-1.9.tar.bz2 && \
    cd htslib-1.9 && ./configure && make && make install

WORKDIR /

# Installing samtools 1.9

RUN wget https://github.com/samtools/samtools/releases/download/1.9/samtools-1.9.tar.bz2 && tar -xjvf samtools-1.9.tar.bz2 && \
    cd samtools-1.9 && ./configure && make && make install

WORKDIR /
