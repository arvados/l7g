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

# Installing SAM tools
RUN apt-get update && apt-get install samtools && apt-get clean
