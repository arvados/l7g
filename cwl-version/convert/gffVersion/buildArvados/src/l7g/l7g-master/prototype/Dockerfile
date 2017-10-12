FROM debian

RUN apt-get update && \
  apt-get install -y git jq wget zip at \
  build-essential curl libcurl4-openssl-dev fuse-emulator-utils \
  libmicrohttpd-dev golang libhts-dev zlib1g-dev vim-common pkg-config \
  supervisor openssh-server \
  --fix-missing && \
  mkdir -p /data/config

RUN cd /root && \
  mkdir libconfig && cd libconfig && \
  wget http://www.hyperrealm.com/libconfig/libconfig-1.5.tar.gz && \
  tar xvfz libconfig-1.5.tar.gz && \
  cd libconfig-1.5 && \
  ./configure && make && make install



# TODO:
#
# * hg19.json for glfd
# * all other auxiliary files
#
#  mkdir -p /home/lightning && \
RUN adduser --disabled-password --gecos "" lightning && \
  cd /home/lightning && \
  mkdir .go && \
  export GOPATH=/home/lightning/.go && \
  go get github.com/robertkrimen/otto && \
  go get github.com/abeconnelly/sloppyjson && \
  go get github.com/abeconnelly/pasta && \
  go get github.com/abeconnelly/pasta/gvcf && \
  go get github.com/aebruno/twobit && \
  go get github.com/julienschmidt/httprouter && \
  go get github.com/mattn/go-sqlite3 && \
  go get github.com/abeconnelly/cglf && \
  go get github.com/abeconnelly/dlug && \
  go get github.com/abeconnelly/autoio && \
  git clone https://github.com/abeconnelly/glfd && \
    cd glfd && cd c/asmukk && make && cd ../.. && ln -s c/asmukk/asm_ukk.a . && go build && cd /home/lightning && \
  git clone https://github.com/abeconnelly/lci && \
    cd lci && go build && cd .. && \
    cd /home/lightning/lci/example && ln -s ../html/css . && ln -s ../html/js . && cd /home/lightning && \
  git clone https://github.com/abeconnelly/l7g-v5t-clinvar && \
    cd l7g-v5t-clinvar && go build && cd /home/lightning && \
  git clone https://github.com/abeconnelly/l7g-p7e-untap && \
    cd l7g-p7e-untap && go build && cd /home/lightning && \
  git clone https://github.com/abeconnelly/cgf && \
    cd /home/lightning && cd cgf && go build && cd cpp && ./cmp.sh && cd muduk && ./cmp.sh && cd /home/lightning && \
  cd /home/lightning && mkdir -p bin && \
    cp cgf/cpp/muduk/bin/muduk bin && \
    cp glfd/glfd bin && \
    cp l7g-p7e-untap/l7g-p7e-untap bin && \
    cp l7g-v5t-clinvar/l7g-v5t-clinvar bin && \
    cp lci/lci bin && \
  cd /home && \
  chown -R lightning:lightning lightning

COPY etc/supervisor/conf.d/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

#COPY etc/supervisor/conf.d/cgf.conf /etc/supervisor/conf.d/cgf.conf
#COPY etc/supervisor/conf.d/glfd.conf /etc/supervisor/conf.d/glfd.conf
#COPY etc/supervisor/conf.d/l7g-p7e-untap.conf /etc/supervisor/conf.d/l7g-p7e-untap.conf
#COPY etc/supervisor/conf.d/l7g-v5t-clinvar.conf /etc/supervisor/conf.d/l7g-v5t-clinvar.conf
#COPY etc/supervisor/conf.d/lci.conf /etc/supervisor/conf.d/lci.conf

COPY config/lci-config.json /data/config/lci-config.json
COPY config/l7g-p7e-config.json /data/config/l7g-p7e-config.json
COPY config/l7g-v5t-config.json /data/config/lci-config.json
COPY config/tile-server-conf.json /data/config/lci-conf.json
COPY config/muduk.cfg /data/config/muduk.cfg


EXPOSE 22 80 8080 8081 8082 8083 8084 8085
CMD ["/usr/bin/supervisord"]
