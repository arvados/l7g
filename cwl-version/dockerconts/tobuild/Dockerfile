FROM arvados/jobs:latest
USER root

# Installing apt-utils
RUN apt-get install -y --no-install-recommends apt-utils

# installing zlib1g-dec
RUN apt-get install zlib1g-dev -y

# installing build essentials
RUN apt-get install build-essential -y

# installing bzip2
RUN apt-get install bzip2 -y

# installing unzip
RUN apt-get install unzip -y

# installing numpy
RUN apt-get install python-numpy -y

# installing Cmake
RUN apt-get install cmake -y 

# clean up
RUN apt-get clean

# install cnpy
RUN wget https://github.com/rogersce/cnpy/archive/master.zip
RUN unzip master.zip -d lib
RUN mv lib/cnpy-master lib/cnpy
RUN cd lib/cnpy && cmake . && make


