FROM arvados/jobs:latest
#FROM debian:9
USER root

# Installing java8
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y software-properties-common
RUN echo oracle-java7-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections
RUN echo "deb http://ppa.launchpad.net/webupd8team/java/ubuntu xenial main" | tee /etc/apt/sources.list.d/webupd8team-java.list
RUN echo "deb-src http://ppa.launchpad.net/webupd8team/java/ubuntu xenial main" | tee -a /etc/apt/sources.list.d/webupd8team-java.list
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys EEA14886
RUN apt-get update
RUN apt-get install -y oracle-java8-installer

# Installing apt-utils
RUN apt-get install -y --no-install-recommends apt-utils

# Installing vcftools
#RUN apt-get update && apt-get install vcftools -y && apt-get clean

# installing zlib1g
RUN apt-get install zlib1g-dev -y
#RUN apt-get install zlib1g -y

# installing libbz2-dev
RUN apt-get install libbz2-dev -y

# installing build essentials
RUN apt-get install build-essential -y

# installing bzip2
RUN apt-get install bzip2 -y

# install unzip
RUN apt-get install unzip -y

# install lib files to build SAMtools
#RUN apt-get update && apt-get install libncurses5-dev -y && apt-get install libbz2-dev -y && apt-get install liblzma-dev -y && apt-get clean

# Installing SAM tools
#RUN wget https://github.com/samtools/samtools/releases/download/1.5/samtools-1.5.tar.bz2
#RUN bzip2 -d samtools-1.5.tar.bz2
#RUN tar xf samtools-1.5.tar
#RUN cd samtools-1.5 && ./configure && make && make install

# Installing parallel
RUN apt-get install parallel -y

# Installing numpy
RUN apt-get install python-numpy -y

# Installing Cmake
RUN apt-get install cmake -y 

# Installing golang

RUN apt-get install golang -y

# Installing xz-utils

RUN apt-get install xz-utils -y

# Installing liblzma-dev
RUN apt-get install liblzma-dev -y

# Installing pthreads
RUN apt-get install libpthread-stubs0-dev -y

# Installing libssl-dev
RUN apt-get install libssl-dev -y

# Installing htslib
#RUN apt-get install libhts-dev -y

# Installing libsdsl

#RUN apt-get install libsdsl-dev

# Clean up
RUN apt-get clean

# Install cnpy
# Using l7g, contains a bugfix not in original
#RUN wget https://github.com/rogersce/cnpy/archive/master.zip
RUN wget https://github.com/curoverse/l7g/archive/master.zip
RUN unzip master.zip -d /lib
RUN rm -r master.zip
RUN mv lib/l7g-master/lib/cnpy lib/cnpy
RUN rm -r lib/l7g-master
RUN cd lib/cnpy && cmake . && make

# Installing htslib
RUN cd /tmp && wget -O- 'https://github.com/samtools/htslib/releases/download/1.6/htslib-1.6.tar.bz2' | tar xjf - && cd htslib-1.6 && make && make install && cd /tmp && rm -rf htslib-1.6

# Installing SDSL
RUN cd /tmp && git clone 'https://github.com/simongog/sdsl-lite.git' && cd sdsl-lite && ./install.sh /usr/local
