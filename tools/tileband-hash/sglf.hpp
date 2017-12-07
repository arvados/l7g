#ifndef SGLF_HPP
#define SGLF_HPP

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <openssl/md5.h>

#include <vector>
#include <string>

typedef struct sglf_type {
  std::vector< std::vector< std::vector< std::string > > > seq;
  std::vector< std::vector< std::vector< int > > > span;
} sglf_t;

int sglf_read(FILE *ifp, sglf_t &sglf);
int sglf_get(std::string &seq, int &span, int tilepath, int tilever, int tilestep, int tilevar, sglf_t &sglf);
int sglf_print(sglf_t &sglf, int tilepath);

#endif
