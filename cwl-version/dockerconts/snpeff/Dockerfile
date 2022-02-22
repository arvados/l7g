FROM ubuntu:18.04
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
        python3 \
        git \
        default-jdk \
        unzip

# Getting HTSLIB 1.10.2 for tabix/bgzip

RUN wget https://github.com/samtools/htslib/releases/download/1.10.2/htslib-1.10.2.tar.bz2 && tar -xjvf htslib-1.10.2.tar.bz2 && \
    cd htslib-1.10.2 && ./configure && make && make install

WORKDIR /

# Installing samtools 1.10

RUN wget https://github.com/samtools/samtools/releases/download/1.10/samtools-1.10.tar.bz2 && tar -xjvf samtools-1.10.tar.bz2 && \
    cd samtools-1.10 && ./configure && make && make install

WORKDIR /

# Installing bcftools 1.10.2

RUN wget https://github.com/samtools/bcftools/releases/download/1.10.2/bcftools-1.10.2.tar.bz2 && tar -xjvf bcftools-1.10.2.tar.bz2 && \
    cd bcftools-1.10.2 && ./configure && make && make install

WORKDIR /

# Getting snpEff4.3t

RUN wget http://sourceforge.net/projects/snpeff/files/snpEff_v4_3t_core.zip && \
    unzip snpEff_v4_3t_core.zip

WORKDIR /
