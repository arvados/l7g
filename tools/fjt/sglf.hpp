#ifndef SGLF_HPP
#define SGLF_HPP

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include <vector>
#include <map>
#include <string>
//#include <pair>

#include "twoBit.hpp"

typedef struct sglf_path_type {
  std::vector< std::vector< std::string > > ext_tileid;
  std::vector< std::vector< std::string > > hash;
  std::vector< std::vector< std::string > > seq;

  std::map< std::string, std::pair< int, int > > hash_pos;
} sglf_path_t;

int sglf_path_step_lookup_hash_variant_id(sglf_path_t *sp, int tilestep, std::string &hash);
int sglf_path_step_lookup_seq_variant_id(sglf_path_t *sp, int tilestep, std::string &seq);

void sglf_path_print(sglf_path_t *sp);

uint16_t tileid_part(uint64_t tileid, int part);
uint64_t parse_tileid(const char *tileid);
int read_sglf_path(FILE *ifp, sglf_path_t &sp);


//--


/*
typedef struct sglf2bit_type {
  twoBit_t twobit;
  int freq;
  sglf2bit_type() : twobit(NULL), freq(0) { }
} sglf2bit_t;
*/

typedef struct sglf2bit_tilepath_type {
  std::vector< std::vector< std::string > > ext_tileid;
  std::vector< std::vector< std::string > > hash;
  std::vector< std::vector< twoBit_t > > seq2bit;

  std::map< std::string, std::pair< int, int > > hash_pos;
} sglf2bit_tilepath_t;

int read_sglf2bit_tilepath(FILE *ifp, sglf2bit_tilepath_t &sp);
int sglf2bit_tilepath_step_lookup_seq_variant_id(sglf2bit_tilepath_t &sp, int tilestep, std::string &seq);


#endif
