#ifndef SGLF_HPP
#define SGLF_HPP

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <openssl/md5.h>

#include <vector>
#include <string>

#include "twoBit.hpp"

enum SGLF_TYPE_ENUM {
  SGLF_TYPE_SEQ = 0,
  SGLF_TYPE_2BIT
} ;

typedef struct sglf_type {
  std::vector< std::vector< std::vector< std::string > > > seq;
  std::vector< std::vector< std::vector< twoBit_t > > > seq2bit;
  std::vector< std::vector< std::vector< int > > > span;

  int type;

  sglf_type() : type(SGLF_TYPE_SEQ) { }

} sglf_t;

int sglf_read(FILE *ifp, sglf_t &sglf);
int sglf_read_2bit(FILE *ifp, sglf_t &sglf);

//int sglf_get(std::string &seq, int &span, int tilepath, int tilever, int tilestep, int tilevar, sglf_t &sglf);

int sglf_print(sglf_t &sglf, int tilepath);

#endif
