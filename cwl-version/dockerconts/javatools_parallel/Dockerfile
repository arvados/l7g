FROM arvados/jobs:latest
USER root

# Installing java8
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y software-properties-common
RUN apt-get update
RUN echo oracle-java7-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections
RUN echo "deb http://ppa.launchpad.net/webupd8team/java/ubuntu xenial main" | tee /etc/apt/sources.list.d/webupd8team-java.list
RUN echo "deb-src http://ppa.launchpad.net/webupd8team/java/ubuntu xenial main" | tee -a /etc/apt/sources.list.d/webupd8team-java.list
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys EEA14886
RUN apt-get update
RUN apt-get install -y oracle-java8-installer

# Installing apt-utils
RUN apt-get update && apt-get install -y --no-install-recommends apt-utils

# Installing vcftools
RUN apt-get update && apt-get install vcftools -y && apt-get clean

# installing zlib1g-dec
RUN apt-get update && apt-get install zlib1g-dev -y && apt-get clean

# installing build essentials
RUN apt-get update && apt-get install build-essential -y && apt-get clean

# installing bzip2
RUN apt-get update && apt-get install bzip2 -y && apt-get clean

# install lib files to build SAMtools
RUN apt-get update && apt-get install libncurses5-dev -y && apt-get install libbz2-dev -y && apt-get install liblzma-dev -y && apt-get clean

# Installing SAM tools
RUN wget https://github.com/samtools/samtools/releases/download/1.5/samtools-1.5.tar.bz2
RUN bzip2 -d samtools-1.5.tar.bz2
RUN tar xf samtools-1.5.tar
RUN cd samtools-1.5 && ./configure && make && make install

# Installing parallel
RUN apt-get update && apt-get install parallel -y && apt-get clean
