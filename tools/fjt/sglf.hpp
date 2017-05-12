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

typedef struct sglf_path_type {
  std::vector< std::vector< std::string > > ext_tileid;
  std::vector< std::vector< std::string > > hash;
  std::vector< std::vector< std::string > > seq;

  std::map< std::string, std::pair< int, int > > hash_pos;
} sglf_path_t;

int sglf_path_step_lookup_hash_variant_id(sglf_path_t *sp, int tilestep, std::string &hash);
int sglf_path_step_lookup_seq_variant_id(sglf_path_t *sp, int tilestep, std::string &seq);

void sglf_path_print(sglf_path_t *sp);

#endif
