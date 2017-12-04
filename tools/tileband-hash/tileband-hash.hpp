#ifndef TILEBAND_HASH_HPP
#define TILEBAND_HASH_HPP

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <openssl/md5.h>

#include <vector>
#include <string>

typedef struct band_info_type {
  std::vector< int > band[2];
  std::vector< std::vector< int > > noc[2];
} band_info_t;

int read_bands(FILE *ifp, std::vector< band_info_t > &band_info_v);
int read_band(FILE *ifp, band_info_t &band_info);
void band_print(band_info_t &band_info);
void print_bands(std::vector< band_info_t > &band_info_v);

int band_md5_hash(std::vector< std::string > &digest_str,
                  std::vector< band_info_t > &band_v,
                  sglf_t &sglf,
                  std::vector< int > &tilepath_list);

#endif
