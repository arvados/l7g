FROM ubuntu:16.04
MAINTAINER Jiayong Li <jli@curii.com>
USER root

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

#Installing OpenJDK-8
RUN apt-get update && \
    apt-get install -y openjdk-8-jdk && \
    apt-get install -y ant && \
    apt-get clean

#Fixing certificate issues
RUN apt-get update && \
    apt-get install ca-certificates-java && \
    apt-get clean && \
    update-ca-certificates -f

# Getting beagle5.4
RUN wget https://faculty.washington.edu/browning/beagle/beagle.05May22.33a.jar
RUN wget https://faculty.washington.edu/browning/beagle/bref3.05May22.33a.jar

# Getting HTSLIB 1.9 for tabix/bgzip

RUN wget https://github.com/samtools/htslib/releases/download/1.9/htslib-1.9.tar.bz2 && tar -xjvf htslib-1.9.tar.bz2 && \
    cd htslib-1.9 && ./configure && make && make install

WORKDIR /

# Installing samtools 1.9

RUN wget https://github.com/samtools/samtools/releases/download/1.9/samtools-1.9.tar.bz2 && tar -xjvf samtools-1.9.tar.bz2 && \
    cd samtools-1.9 && ./configure && make && make install

WORKDIR /

# Installing bcftools 1.9

RUN wget https://github.com/samtools/bcftools/releases/download/1.9/bcftools-1.9.tar.bz2 && tar -xjvf bcftools-1.9.tar.bz2 && \
    cd bcftools-1.9 && ./configure && make && make install

WORKDIR /
