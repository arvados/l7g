#ifndef CGFT_H
#define CGFT_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <unistd.h>
#include <string.h>
#include <math.h>
#include <errno.h>
#include <sys/time.h>

#include <getopt.h>

#include <cstdlib>
#include <map>
#include <vector>
#include <string>
#include <complex>
#include <iterator>


#include <sdsl/vectors.hpp>
#include <sdsl/bit_vectors.hpp>

#define CGFT_MAGIC "{\"cgf.b\""
#define CGF_VERSION "0.3.1"
#define CGLF_VERSION "0.1.0"

#define OVF16_MAX 0xffff
#define OVF64_MAX 0xffffffffffffffff

#define SPAN_SDSL_ENC_VAL (1<<30)

extern char DEFAULT_TILEMAP[];

typedef struct tilepath_type {
  uint64_t TilePath;
  std::string Name;
  uint64_t NTileStep;
  uint64_t NOverflow;
  uint64_t NOverflow64;
  uint64_t ExtraDataSize;

  unsigned char *Loq;
  unsigned char *Span;
  uint64_t *Cache;

  uint16_t *Overflow;
  uint64_t *Overflow64;

  char *ExtraData;

  // Low quality information.
  // The "Het" and "Hom" portions indicate whether
  // the low quality data is same across both alleles
  // and doesn't indicate anything about whether the
  // variants are heterozygous or homozygous.
  // Hom represents the bulk of the data so
  // we get size savings by splitting it out
  // into Het and Hom parts.
  //
  // - Step is strictly increasing TileStep positions
  // - Variant are allele interleaved tile variants
  // - NocSum is the current inclusive count of the Noc elements
  // - NocStart is the start of the nocall run
  // - NocLen is the run of nocall elements
  //
  // Since the TileStep and NocSum are strictly non-decreasing,
  // they benefit from being an 'enc_vector' whereas the rest
  // benefit from being variable length encoded.
  //

  uint64_t LoqTileStepHomSize;
  uint64_t LoqTileVariantHomSize;
  uint64_t LoqTileNocSumHomSize;
  uint64_t LoqTileNocStartHomSize;
  uint64_t LoqTileNocLenHomSize;

  uint64_t LoqTileStepHetSize;
  uint64_t LoqTileVariantHetSize;
  uint64_t LoqTileNocSumHetSize;
  uint64_t LoqTileNocStartHetSize;
  uint64_t LoqTileNocLenHetSize;

  sdsl::enc_vector<> LoqTileStepHom;
  sdsl::vlc_vector<> LoqTileVariantHom;
  sdsl::enc_vector<> LoqTileNocSumHom;
  sdsl::vlc_vector<> LoqTileNocStartHom;
  sdsl::vlc_vector<> LoqTileNocLenHom;

  sdsl::enc_vector<> LoqTileStepHet;
  sdsl::vlc_vector<> LoqTileVariantHet;
  sdsl::enc_vector<> LoqTileNocSumHet;
  sdsl::vlc_vector<> LoqTileNocStartHet;
  sdsl::vlc_vector<> LoqTileNocLenHet;

} tilepath_t;

typedef struct cgf_type {
  unsigned char Magic[8];
  std::string CGFVersion;
  std::string LibraryVersion;
  uint64_t PathCount;
  std::string TileMap;
  std::vector<uint64_t> PathStructOffset;
  std::vector<tilepath_t> Path;
} cgf_t;

//----

typedef struct tilepath_vec_type {
  std::string name;
  std::vector<int> allele[2];
  std::vector<int> loq_flag[2];
  std::vector< std::vector<int> > loq_info[2];
} tilepath_vec_t;

typedef struct tilepath_ez_type {

  int tilepath;

  int N;

  // hiq
  //
  std::vector<uint64_t> cache;
  std::vector<unsigned char> span_bv;

  // interleaved overflow
  //
  std::vector<int16_t> ovf_vec;
  std::vector<int32_t> ovf32_vec;
  std::vector<int64_t> ovf64_vec;

  std::vector<char> data_vec;

  // loq
  //
  std::vector<unsigned char> loq_bv;  // floor( (N+7)/8 )

  int n_loq;

  std::vector<int> loq_info_pos;

  // interleaved for multi-allelic
  //
  std::vector<int> loq_info_variant;
  std::vector<int> loq_info_sn;
  std::vector<int> loq_info_noc;

  std::vector<int> loq_info_pos_hom;
  std::vector<int> loq_info_variant_hom;
  std::vector<int> loq_info_sn_hom;
  std::vector<int> loq_info_noc_hom;

  std::vector<int> loq_info_pos_het;
  std::vector<int> loq_info_variant_het;
  std::vector<int> loq_info_sn_het;
  std::vector<int> loq_info_noc_het;

} tilepath_ez_t;

typedef struct cgf_ez_type {

  std::string tilemap_str;

  std::map< std::string, int > tilemap;
  std::vector<tilepath_ez_t> tilepath;

} cgf_ez_t;

int ez_save(const char *base_f, int tilepath, tilepath_ez_t &ez);

const char *read_tilemap_from_file(std::string &, const char *);

//void cgft_create_container(FILE *, const char *);
void cgft_create_container(FILE *, const char *, const char *, const char *);
void cgft_print_header(cgf_t *);
void cgft_print_tilepath(cgf_t *, tilepath_t *);

cgf_t *cgft_read(FILE *);

void cgft_tilepath_init(tilepath_t &, uint64_t);

int cgft_read_band_tilepath(cgf_t *, tilepath_t *, FILE *);
int cgft_sanity(cgf_t *);

int cgft_write_to_file(cgf_t *, const char *);

int cgft_output_band_format(cgf_t *, tilepath_t *, FILE *);

// ez functions

void print_bgf(tilepath_vec_t &);
void print_tilepath_vec(tilepath_vec_t &);
void ez_create(tilepath_ez_t &, tilepath_vec_t &, std::map< std::string, int > &);
void ez_print(tilepath_ez_t &);

int load_tilemap(std::string &, std::map< std::string, int > &);
void mk_tilemap_key(std::string &key, tilepath_vec_t &tilepath, int tilestep, int n);

// helper functions

inline int NumberOfSetBits32(uint32_t u)
{
  u = u - ((u >> 1) & 0x55555555);
  u = (u & 0x33333333) + ((u >> 2) & 0x33333333);
  return (((u + (u >> 4)) & 0x0F0F0F0F) * 0x01010101) >> 24;
}

// This is slower than the above but is more explicit
//
inline int NumberOfSetBits8(uint8_t u)
{
  u = (u & 0x55) + ((u>>1) & 0x55);
  u = (u & 0x33) + ((u>>2) & 0x33);
  u = (u & 0x0f) + ((u>>4) & 0x0f);
  return u;
}



#endif

