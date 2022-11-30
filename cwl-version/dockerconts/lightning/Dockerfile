FROM ubuntu:latest
MAINTAINER Jiayong Li <jli@curii.com>
USER root
ARG DEBIAN_FRONTEND=noninteractive

# Install necessary dependencies

RUN apt-get update --fix-missing -qy

RUN apt-get install -qy --no-install-recommends wget \
  build-essential \
  cmake \
  zlib1g-dev \
  libbz2-dev \
  liblzma-dev \
  libncurses5-dev \
  libncursesw5-dev \
  gcc \
  make \
  python3.8 \
  python3-pip \
  python3-numpy \
  python3-pandas \
  python3-scipy \
  python3-matplotlib \
  git

RUN pip3 install sklearn
RUN pip3 install --upgrade scipy

# Installing go 1.19

RUN wget https://go.dev/dl/go1.19.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin:/root/go/bin

# Getting lightning

RUN git clone https://github.com/arvados/lightning && \
    cd lightning && make

WORKDIR /
