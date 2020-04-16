FROM ubuntu:latest
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
        vcftools \
        gcc \
        libperl-dev \
        dtrx \
        autoconf \
        make \
        libpq-dev \
        python \
        python-pip \
        python3 \
        python3-pip \
        libmysqlclient-dev \
        git \
        curl \
        g++ \
        libgsl0-dev \
        swig \
        perl \
        default-jdk

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

RUN pip install biopython

RUN git clone https://github.com/mutalyzer/crossmapper && \
    cd crossmapper && python setup.py install

WORKDIR /

# Installing mutalyzer description-extractor

RUN pip install description-extractor

# Install HGVS

RUN pip3 install --upgrade setuptools && \
    pip3 install Cython && pip3 install hgvs

# Installing counsyl

RUN git clone https://github.com/counsyl/hgvs.git && \
    cd hgvs && python setup.py install

WORKDIR /

# Getting snpEff

RUN wget http://sourceforge.net/projects/snpeff/files/snpEff_v4_3t_core.zip && \
    unzip snpEff_v4_3t_core.zip
