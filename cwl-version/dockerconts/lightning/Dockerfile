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

# Installing go 1.16

RUN wget https://golang.org/dl/go1.16.7.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.16.7.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin:/root/go/bin

# Getting lightning

RUN git clone https://github.com/arvados/lightning && \
    cd lightning && make

WORKDIR /

