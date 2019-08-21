/*
    Copyright Curoverse, Inc.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// FastJ Tool
//
// An attempt at a tool to do FastJ manipulation
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <getopt.h>

#include <openssl/md5.h>
#include <zlib.h>

#include <vector>
#include <string>
#include <map>
#include <algorithm>

#include "cJSON.h"
#include "sglf.hpp"
//#include "sglf2bit.hpp"

//#define FASTJ_TOOL_DEBUG 1

#define FASTJ_TOOL_VERSION "0.1.5"

// structures of this type will be populated with the info about each tile
//
typedef struct fj_tile_type {
  cJSON *hdr;
  std::string seq;
  uint64_t tileid;
  int span;
  int n;
  int nocall;
  std::string startTag, endTag;
  std::string seqHash;
  int startTile;
  int endTile;
} fj_tile_t;

// these are used to indicate which command line option was specified
//
enum FJT_ACTION {
 FJT_NOOP = 0,
 FJT_CSV,
 FJT_CONCAT,
 FJT_FILTER,
 FJT_BAND,
 FJT_BAND_CONVERT,
 FJT_BAND_BATCH_HASH,
 FJT_TEST,
};

typedef struct band_info_type {
  std::vector< int > band[2];
  std::vector< std::vector< int > > noc[2];
} band_info_t;

int verbose_flag = 0;

static struct option long_options[] = {
  {"help", no_argument, NULL, 'h'},
  {"verbose", no_argument, NULL, 'v'},
  {"version", no_argument, NULL, 'V'},
  {"csv", no_argument, NULL, 'C'},
  {"band", no_argument, NULL, 'B'},
  {"band-convert", no_argument, NULL, 'b'},
  {"band-batch-hash", no_argument, NULL, 'H'},
  {"concatenate", required_argument, NULL, 'c'},
  {"tile-path", required_argument, NULL, 'p'},
  {"tile-library", required_argument, NULL, 'L'},
  {"input", required_argument, NULL, 'i'},
  {"input-file-list", required_argument, NULL, 'I'},
  {"unsorted", no_argument, NULL, 'U'},
  {"test", no_argument, NULL, 'T'},
  {"test-tileid", no_argument, NULL, 't'},
  {0,0,0,0}
};

void show_version() {
  printf("fjt version: %s\n", FASTJ_TOOL_VERSION);
}

void show_help() {
  show_version();
  printf("usage:\n  fjt [-c variant] [-C] [-v] [-V] [-h] [input]\n");
  printf("\n");
  printf("  [-C]            Output comma separated `extended tileID`, `hash` and `sequence` (CSV output)\n");
  printf("  [-B]            Output band format\n");
  printf("  [-b]            input band format and output FastJ (requires '-L sglf' option)\n");
  printf("  [-H]            batch hash of input bands (requires '-L sglf' option)\n");
  printf("  [-c variant]    Concatenate FastJ tiles into sequence.  `variant` is the variant ID to concatenate on\n");
  printf("  [-L sglf]       Simple genome library format tile path file\n");
  printf("  [-i ifn]        input file\n");
  printf("  [-I ifn_list]   file containing gziped list of FastJ files to convert\n");
  printf("  [-p tilepath]   Tile path (in decimal)\n");
  printf("  [-U]            do not sort output (for use with -C option)\n");
  printf("  [-T]            Check fastj file for integrity.\n");
  printf("  [-t]            Check fastj tileIDs for consistency.\n");
  printf("  [-v]            Verbose\n");
  printf("  [-V]            Version\n");
  printf("  [-h]            Help\n");
  printf("\n");
}

void md5str(std::string &s, std::string &seq) {
  int i;
  unsigned char m[MD5_DIGEST_LENGTH];
  char buf[32];

  s.clear();

  MD5((unsigned char *)(seq.c_str()), seq.size(), m);

  for (i=0; i<MD5_DIGEST_LENGTH; i++) {
    sprintf(buf, "%02x", (unsigned char)m[i]);
    s += buf;
  }
}


bool sortTileCmp(const fj_tile_t &lhs, const fj_tile_t &rhs) {
  return lhs.tileid < rhs.tileid;
}

enum fj_input_state { EXPECT_HDR, READ_HDR, READ_SEQ };

void print_tileid(uint64_t tileid) {
  uint64_t u64;
  int curpos=0;
  unsigned int byte_offset[] = { 6, 4, 2, 0 };
  std::string ofmt[] = { "%04x", "%02x", "%04x", "%03x" };

  for (curpos=0; curpos<4; curpos++) {
    u64 = tileid>>(8*byte_offset[curpos]);
    u64 &= 0xffff;
    if (curpos>0) { printf("."); }
    printf( ofmt[curpos].c_str(), (unsigned int)u64 );
  }
}

int read_bands(FILE *ifp, std::vector< band_info_t > &band_info_v) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;
  std::vector< int > noc_vec;

  std::string buf;

  int read_state = 0;
  int bracket_count=0;
  int cur_val=-3;

  band_info_t band_info;

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }

    char_no++;
    if (ch=='\n') {
      line_no++;

      switch(read_state) {
        case 0:
          read_state++;
          break;
        case 1:
          read_state++;
          break;
        case 2:
          read_state++;
          break;
        case 3:
          read_state=0;
          bracket_count=0;
          buf.clear();
          band_info_v.push_back(band_info);
          band_info.band[0].clear();
          band_info.band[1].clear();
          band_info.noc[0].clear();
          band_info.noc[1].clear();
          break;
        default:
          return -1;
      }
      continue;
    }

    if (ch==' ') {
      if (buf.size()>0) {
        cur_val = atoi(buf.c_str());

        if (read_state < 2) {
          band_info.band[read_state].push_back(cur_val);
        }

        else {
          noc_vec.push_back(cur_val);
        }

      }
      buf.clear();
      continue;
    }

    if (ch=='[') { bracket_count++; continue; }
    if (ch==']') {
      bracket_count--;

      // Tile variant bands still
      //
      if (read_state<2) {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());

          if (read_state < 2) {
            band_info.band[read_state].push_back(cur_val);
          }
          buf.clear();
        }

      }

      // nocall information
      //
      else {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());
          noc_vec.push_back(cur_val);
          buf.clear();
        }

        if (bracket_count==1) {
          band_info.noc[read_state-2].push_back(noc_vec);
          noc_vec.clear();
        }
      }

      continue;
    }

    buf += (char)ch;

  }

  return 0;
}

int read_band(FILE *ifp, band_info_t &band_info) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;
  std::vector< int > noc_vec;

  std::string buf;

  int read_state = 0;
  int bracket_count=0;
  int cur_val=-3;

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }
    char_no++;
    if (ch=='\n') {
      line_no++;

      switch(read_state) {
        case 0:
          break;
        case 1:
          break;
        case 2:
          break;
        case 3:
          break;
        default:
          return -1;
      }
      read_state++;
      continue;
    }

    if (ch==' ') {
      if (buf.size()>0) {
        cur_val = atoi(buf.c_str());

        if (read_state < 2) {
          band_info.band[read_state].push_back(cur_val);
        }

        else {
          noc_vec.push_back(cur_val);
        }

      }
      buf.clear();
      continue;
    }

    if (ch=='[') { bracket_count++; continue; }
    if (ch==']') {
      bracket_count--;

      // Tile variant bands still
      //
      if (read_state<2) {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());

          if (read_state < 2) {
            band_info.band[read_state].push_back(cur_val);
          }
          buf.clear();
        }

      }

      // nocall information
      //
      else {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());
          noc_vec.push_back(cur_val);
          buf.clear();
        }

        if (bracket_count==1) {
          band_info.noc[read_state-2].push_back(noc_vec);
          noc_vec.clear();
        }
      }

      continue;
    }

    buf += (char)ch;

  }

  return 0;
}


int read_tiles_and_print_csv(FILE *ifp) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;

  std::string m5;

  fj_tile_t cur_tile;

  fj_input_state state;
  std::string buf;

  state = EXPECT_HDR;

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }
    char_no++;
    if (ch=='\n') { line_no++; }

    if (state==EXPECT_HDR) {
      if (ch=='\n') { continue; }
      if (ch==' ') { continue; }

      if (ch=='>') {

        // add to list
        //
        if (buf.size()>0) {
          cur_tile.seq = buf;

          if (cur_tile.hdr==NULL) { return -8; }
          cJSON *tid = cjson_obj(cur_tile.hdr, "tileID");
          if (cJSON_IsString(tid)) {
            cur_tile.tileid = parse_tileid(tid->valuestring);

            cJSON *span = cjson_obj(cur_tile.hdr, "seedTileLength");
            if (cJSON_IsNumber(span)) {
              cur_tile.span = (int)(span->valuedouble);
            } else {
              return -1;
            }
          } else { return -4; }

          print_tileid(cur_tile.tileid);
          printf("+%x", cur_tile.span);
          md5str(m5, cur_tile.seq);
          printf(",%s", m5.c_str());
          printf(",%s\n", cur_tile.seq.c_str());

          //fj_tile.push_back(cur_tile);
          if (cur_tile.hdr) { cJSON_Delete(cur_tile.hdr); }
          cur_tile.hdr = NULL;
          cur_tile.seq.clear();
        }

        buf.clear();
        state = READ_HDR;
        continue;
      }

      return -2;
    }

    if (state==READ_HDR) {
      if (ch=='\n') {
        cur_tile.hdr = cJSON_Parse(buf.c_str());
        if (cur_tile.hdr==NULL) { return -3; }
        buf.clear();
        state = READ_SEQ;

        continue;
      }
      buf += (char)ch;
      continue;
    }

    if (state==READ_SEQ) {
      if ((ch==' ') || (ch=='\n')) { continue; }
      if (ch=='>') {

        // add to list
        //
        if (buf.size()>0) {
          cur_tile.seq = buf;

          if (cur_tile.hdr==NULL) { return -9; }
          cJSON *tid = cjson_obj(cur_tile.hdr, "tileID");
          if (cJSON_IsString(tid)) {
            cur_tile.tileid = parse_tileid(tid->valuestring);

            cJSON *span = cjson_obj(cur_tile.hdr, "seedTileLength");
            if (cJSON_IsNumber(span)) {
              cur_tile.span = (int)(span->valuedouble);
            } else {
              return -1;
            }
          } else { return -4; }


          print_tileid(cur_tile.tileid);
          printf("+%x", cur_tile.span);
          md5str(m5, cur_tile.seq);
          printf(",%s", m5.c_str());
          printf(",%s\n", cur_tile.seq.c_str());

          //fj_tile.push_back(cur_tile);
          if (cur_tile.hdr) { cJSON_Delete(cur_tile.hdr); }
          cur_tile.hdr = NULL;

          cur_tile.seq.clear();
        }

        buf.clear();
        state = READ_HDR;
        continue;

      }

      buf += (char)ch;
      continue;

    }

  }

  // add final element to list
  //
  if (buf.size()>0) {
    cur_tile.seq = buf;

    if (cur_tile.hdr==NULL) { return -10; }
    cJSON *tid = cjson_obj(cur_tile.hdr, "tileID");
    if (cJSON_IsString(tid)) {
      cur_tile.tileid = parse_tileid(tid->valuestring);

      cJSON *span = cjson_obj(cur_tile.hdr, "seedTileLength");
      if (cJSON_IsNumber(span)) {
        cur_tile.span = (int)(span->valuedouble);
      } else {
        return -1;
      }
    } else { return -4; }

    print_tileid(cur_tile.tileid);
    printf("+%x", cur_tile.span);
    md5str(m5, cur_tile.seq);
    printf(",%s", m5.c_str());
    printf(",%s\n", cur_tile.seq.c_str());

    if (cur_tile.hdr) { cJSON_Delete(cur_tile.hdr); }
    cur_tile.hdr = NULL;
    cur_tile.seq.clear();
  }

  return 0;
}

int read_tiles_gz(gzFile gz_fp, std::vector< fj_tile_t > &fj_tile) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;
  //std::vector< fj_tile_t > fj_tile;
  fj_tile_t cur_tile;

  fj_input_state state;
  std::string buf;

  state = EXPECT_HDR;

  while (gzeof(gz_fp)==0) {
    ch = gzgetc(gz_fp);
    if (ch<0) { continue; }
    char_no++;
    if (ch=='\n') { line_no++; }

    if (state==EXPECT_HDR) {
      if (ch=='\n') { continue; }
      if (ch==' ') { continue; }

      if (ch=='>') {

        // add to list
        //
        if (buf.size()>0) {
          cur_tile.seq = buf;
          fj_tile.push_back(cur_tile);
          cur_tile.hdr = NULL;
          cur_tile.seq.clear();
        }

        buf.clear();
        state = READ_HDR;
        continue;
      }

      return -2;
    }

    if (state==READ_HDR) {
      if (ch=='\n') {
        cur_tile.hdr = cJSON_Parse(buf.c_str());
        if (cur_tile.hdr==NULL) { return -3; }
        buf.clear();
        state = READ_SEQ;
        continue;
      }
      buf += (char)ch;
      continue;
    }

    if (state==READ_SEQ) {
      if ((ch==' ') || (ch=='\n')) { continue; }
      if (ch=='>') {

        // add to list
        //
        if (buf.size()>0) {
          cur_tile.seq = buf;
          fj_tile.push_back(cur_tile);
          cur_tile.hdr = NULL;
          cur_tile.seq.clear();
        }

        buf.clear();
        state = READ_HDR;
        continue;

      }

      buf += (char)ch;
      continue;

    }

  }

  // add final element to list
  //
  if (buf.size()>0) {
    cur_tile.seq = buf;
    fj_tile.push_back(cur_tile);
    cur_tile.hdr = NULL;
    cur_tile.seq.clear();
  }

  for (i=0; i<fj_tile.size(); i++) {
    cJSON *tid = cjson_obj(fj_tile[i].hdr, "tileID");
    if (cJSON_IsString(tid)) {
      fj_tile[i].tileid = parse_tileid(tid->valuestring);

      cJSON *span = cjson_obj(fj_tile[i].hdr, "seedTileLength");
      if (cJSON_IsNumber(span)) {
        fj_tile[i].span = (int)(span->valuedouble);
      } else {
        return -1;
      }
    } else {
      //printf("ERROR, element %i does not have tileID\n", i);
      return -4;
    }
  }

  std::sort( fj_tile.begin(), fj_tile.end(), sortTileCmp );
}


// "reading" means populating fj_tile_t struct with json header and sequence
//
int read_tiles(FILE *ifp, std::vector< fj_tile_t > &fj_tile) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;
  //std::vector< fj_tile_t > fj_tile;
  fj_tile_t cur_tile;

  fj_input_state state;
  std::string buf;

  state = EXPECT_HDR;

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }
    char_no++;
    if (ch=='\n') { line_no++; }

    if (state==EXPECT_HDR) {
      if (ch=='\n') { continue; }
      if (ch==' ') { continue; }

      if (ch=='>') {

        // add to list
        //
        if (buf.size()>0) {
          cur_tile.seq = buf;
          fj_tile.push_back(cur_tile);
          cur_tile.hdr = NULL;
          cur_tile.seq.clear();
        }

        buf.clear();
        state = READ_HDR;
        continue;
      }

      return -2;
    }

    if (state==READ_HDR) {
      if (ch=='\n') {
        cur_tile.hdr = cJSON_Parse(buf.c_str());
        if (cur_tile.hdr==NULL) { return -3; }
        buf.clear();
        state = READ_SEQ;
        continue;
      }
      buf += (char)ch;
      continue;
    }

    if (state==READ_SEQ) {
      if ((ch==' ') || (ch=='\n')) { continue; }
      if (ch=='>') {

        // add to list
        //
        if (buf.size()>0) {
          cur_tile.seq = buf;
          fj_tile.push_back(cur_tile);
          cur_tile.hdr = NULL;
          cur_tile.seq.clear();
        }

        buf.clear();
        state = READ_HDR;
        continue;

      }

      buf += (char)ch;
      continue;

    }

  }

  // add final element to list
  //
  if (buf.size()>0) {
    cur_tile.seq = buf;
    fj_tile.push_back(cur_tile);
    cur_tile.hdr = NULL;
    cur_tile.seq.clear();
  }

  for (i=0; i<fj_tile.size(); i++) {
    cJSON *tid = cjson_obj(fj_tile[i].hdr, "tileID");
    if (cJSON_IsString(tid)) {
      fj_tile[i].tileid = parse_tileid(tid->valuestring);

      cJSON *span = cjson_obj(fj_tile[i].hdr, "seedTileLength");
      if (cJSON_IsNumber(span)) {
        fj_tile[i].span = (int)(span->valuedouble);
      } else {
        return -1;
      }

      cJSON *n = cjson_obj(fj_tile[i].hdr, "n");
      if (cJSON_IsNumber(n)) {
        fj_tile[i].n = (int)(n->valuedouble);
      } else {
        return -1;
      }

      cJSON *nocall = cjson_obj(fj_tile[i].hdr, "nocallCount");
      if (cJSON_IsNumber(nocall)) {
        fj_tile[i].nocall = (int)(nocall->valuedouble);
      } else {
        return -1;
      }

      cJSON *startTag = cjson_obj(fj_tile[i].hdr, "startTag");
      if (cJSON_IsString(startTag)) {
        fj_tile[i].startTag = startTag->valuestring;
      } else {
        return -1;
      }

      cJSON *endTag = cjson_obj(fj_tile[i].hdr, "endTag");
      if (cJSON_IsString(endTag)) {
        fj_tile[i].endTag = endTag->valuestring;
      } else {
        return -1;
      }

      cJSON *seqHash = cjson_obj(fj_tile[i].hdr, "md5sum");
      if (cJSON_IsString(seqHash)) {
        fj_tile[i].seqHash = seqHash->valuestring;
      } else {
        return -1;
      }

      cJSON *startTile = cjson_obj(fj_tile[i].hdr, "startTile");
      if (cJSON_IsBool(startTile)) {
        if ( cJSON_IsTrue(startTile) ) {
          fj_tile[i].startTile = 1;
        }
        else {
          fj_tile[i].startTile = 0;
        }
      } else {
        return -1;
      }

      cJSON *endTile = cjson_obj(fj_tile[i].hdr, "endTile");
      if (cJSON_IsBool(endTile)) {
        if ( cJSON_IsTrue(endTile) ) {
          fj_tile[i].endTile = 1;
        }
        else {
          fj_tile[i].endTile = 0;
        }
      } else {
        return -1;
      }

    } else {
      //printf("ERROR, element %i does not have tileID\n", i);
      return -4;
    }
  }

  std::sort( fj_tile.begin(), fj_tile.end(), sortTileCmp );
}

void concatenate_tiles(std::vector< fj_tile_t > &fj_tile, uint16_t variantid, std::string &seq) {
  uint16_t vid;
  int i, j, k;
  int offset=0;

  seq.clear();

  for (i=0; i<fj_tile.size(); i++) {
    vid = tileid_part(fj_tile[i].tileid, 0);
    if (vid!=variantid) { continue; }

    seq += (fj_tile[i].seq.c_str() + offset);
    offset=24;
  }
}

void cleanup_tiles(std::vector< fj_tile_t > &fj_tile) {
  int i;
  for (i=0; i<fj_tile.size(); i++) { cJSON_Delete(fj_tile[i].hdr); }
}

int create_band_info_sglf2bit(band_info_t &band_info, std::vector< fj_tile_t > &fj_tile, sglf2bit_tilepath_t &sglf2bit_tilepath) {
  int i, j, fj_idx, n, m;
  int tilestep=0, tilevar=0;
  int prev_tilestep=-1, prev_tilevar=-1;
  int sglf_tilevar=-1;
  std::vector< int > noc_v;

  int noc_start, noc_len;

  band_info.band[0].clear();
  band_info.band[1].clear();

  band_info.noc[0].clear();
  band_info.noc[1].clear();

  n = (int)fj_tile.size();

  for (fj_idx=0; fj_idx<n; fj_idx++) {

    noc_v.clear();

    tilestep = (int)tileid_part(fj_tile[fj_idx].tileid, 1);
    tilevar = (int)tileid_part(fj_tile[fj_idx].tileid, 0);

    if (tilevar>=2) { return -1; }
    if (tilevar<0) { return -2; }
    if (tilestep<0) { return -3; }

    // fill in the non-spanning tile vector positions
    //
    for (i=band_info.band[tilevar].size(); i<tilestep; i++) {
      band_info.band[tilevar].push_back( -1 );
      band_info.noc[tilevar].push_back(noc_v);
    }

    // Add the tile variant to the appropriate band allele
    //
    sglf_tilevar = sglf2bit_tilepath_step_lookup_seq_variant_id(sglf2bit_tilepath, tilestep, fj_tile[fj_idx].seq);
    if (sglf_tilevar<0) {

      fprintf(stderr, "tilestep %i (0x%x) was not found in pos index %i\n", tilestep, tilestep, fj_idx);
      fprintf(stderr, "seq: %s\n", fj_tile[fj_idx].seq.c_str());
      return -4;
    }

    band_info.band[tilevar].push_back(sglf_tilevar);

    // add to nocall band vector
    //
    noc_v.clear();

    m = fj_tile[fj_idx].seq.size();
    noc_start = -1;
    noc_len = 0;
    for (i=0; i<m; i++) {
      if ((fj_tile[fj_idx].seq[i]=='n') ||
          (fj_tile[fj_idx].seq[i]=='N')) {
        if (noc_start>=0) { noc_len++; }
        else { noc_start = i; noc_len=1; }
      }

      else {
        if (noc_start>=0) {
          noc_v.push_back(noc_start);
          noc_v.push_back(noc_len);
        }
        noc_start=-1;
        noc_len=0;
      }
    }

    if (noc_start>=0) {
      noc_v.push_back(noc_start);
      noc_v.push_back(noc_len);
    }

    band_info.noc[tilevar].push_back(noc_v);
    int zz = (int)band_info.noc[tilevar].size();

  }


  // special case at the ned if we have spanning tiles.
  // If we do, fill in the remaining tiles with -1
  //
  fj_idx = (int)(fj_tile.size()-1);

  // our final tilestep position
  //

  n = (int)tileid_part(fj_tile[fj_idx].tileid, 1);
  n += fj_tile[fj_idx].span;

  noc_v.clear();
  for (i=band_info.band[0].size(); i<n; i++) { band_info.band[0].push_back(-1); }
  for (i=band_info.band[1].size(); i<n; i++) { band_info.band[1].push_back(-1); }

  for (i=band_info.noc[0].size(); i<n; i++) { band_info.noc[0].push_back(noc_v); }
  for (i=band_info.noc[1].size(); i<n; i++) { band_info.noc[1].push_back(noc_v); }

  return 0;
}



int create_band_info(band_info_t &band_info, std::vector< fj_tile_t > &fj_tile, sglf_path_t &sglf_path) {
  int i, j, fj_idx, n, m;
  int tilestep=0, tilevar=0;
  int prev_tilestep=-1, prev_tilevar=-1;
  int sglf_tilevar=-1;
  //std::vector< int > band[2];
  //std::vector< std::vector< int > > noc_band[2];
  std::vector< int > noc_v;

  int noc_start, noc_len;


  band_info.band[0].clear();
  band_info.band[1].clear();

  band_info.noc[0].clear();
  band_info.noc[1].clear();

  n = (int)fj_tile.size();

  for (fj_idx=0; fj_idx<n; fj_idx++) {

    noc_v.clear();

    tilestep = (int)tileid_part(fj_tile[fj_idx].tileid, 1);
    tilevar = (int)tileid_part(fj_tile[fj_idx].tileid, 0);

    if (tilevar>=2) { return -1; }
    if (tilevar<0) { return -2; }
    if (tilestep<0) { return -3; }

    // fill in the non-spanning tile vector positions
    //
    for (i=band_info.band[tilevar].size(); i<tilestep; i++) {
      band_info.band[tilevar].push_back( -1 );
      band_info.noc[tilevar].push_back(noc_v);
    }


    // Add the tile variant to the appropriate band allele
    //
    sglf_tilevar = sglf_path_step_lookup_seq_variant_id(&sglf_path, tilestep, fj_tile[fj_idx].seq);
    if (sglf_tilevar<0) {

      fprintf(stderr, "tilestep %i (0x%x) was not found in pos index %i\n", tilestep, tilestep, fj_idx);
      fprintf(stderr, "seq: %s\n", fj_tile[fj_idx].seq.c_str());
      return -4;
    }

    band_info.band[tilevar].push_back(sglf_tilevar);

    // add to nocall band vector
    //
    noc_v.clear();

    m = fj_tile[fj_idx].seq.size();
    noc_start = -1;
    noc_len = 0;
    for (i=0; i<m; i++) {
      if ((fj_tile[fj_idx].seq[i]=='n') ||
          (fj_tile[fj_idx].seq[i]=='N')) {
        if (noc_start>=0) { noc_len++; }
        else { noc_start = i; noc_len=1; }
      }

      else {
        if (noc_start>=0) {
          noc_v.push_back(noc_start);
          noc_v.push_back(noc_len);
        }
        noc_start=-1;
        noc_len=0;
      }
    }

    if (noc_start>=0) {
      noc_v.push_back(noc_start);
      noc_v.push_back(noc_len);
    }

    band_info.noc[tilevar].push_back(noc_v);
    int zz = (int)band_info.noc[tilevar].size();

  }


  // special case at the ned if we have spanning tiles.
  // If we do, fill in the remaining tiles with -1
  //
  fj_idx = (int)(fj_tile.size()-1);

  // our final tilestep position
  //

  n = (int)tileid_part(fj_tile[fj_idx].tileid, 1);
  n += fj_tile[fj_idx].span;

  noc_v.clear();
  for (i=band_info.band[0].size(); i<n; i++) { band_info.band[0].push_back(-1); }
  for (i=band_info.band[1].size(); i<n; i++) { band_info.band[1].push_back(-1); }

  for (i=band_info.noc[0].size(); i<n; i++) { band_info.noc[0].push_back(noc_v); }
  for (i=band_info.noc[1].size(); i<n; i++) { band_info.noc[1].push_back(noc_v); }

  return 0;
}

void band_print(band_info_t &band_info) {
  int i, j, a;

  for (a=0; a<2; a++) {
    printf("[");
    for (i=0; i<band_info.band[a].size(); i++) {
      printf(" %i", band_info.band[a][i]);
    }
    printf("]\n");
  }

  for (a=0; a<2; a++) {
    printf("[");
    for (i=0; i<band_info.noc[a].size(); i++) {
      printf("[");
      for (j=0; j<band_info.noc[a][i].size(); j++) {
        printf(" %i", band_info.noc[a][i][j]);
      }
      printf(" ]");
    }
    printf("]\n");
  }
}

void print_bands(std::vector< band_info_t > &band_info_v) {
  int ii;
  for (ii=0; ii<band_info_v.size(); ii++) {
    band_print(band_info_v[ii]);
  }
}


void print_fold(std::string &s, int fold_w) {
  int pos=0;
  for (pos=0; pos<s.size(); pos++) {
    if ((pos>0) && ((pos%fold_w)==0)) { printf("\n"); }
    printf("%c", s[pos]);
  }
  printf("\n");
}

void print_substr(std::string &s, int beg, int n) {
  int i, m;
  m = ( ((int)s.size() < (beg+n)) ? ((int)s.size() - beg) : n );
  for (i=beg; i<(beg+m); i++) { printf("%c", s[i]); }
}

int band_hash(std::string &m5_s, band_info_t &band_info, sglf_path_t &sglf_path) {
  int i, j, k;
  int n, m, a;

  int allele=0;
  int tilestep=0, tilevar=0, span_len=0;
  int noc_count=0;
  int noc_start, noc_len, pos;
  int fold_w = 50;

  std::string mm;

  char *chp;

  MD5_CTX md5_ctx;
  unsigned char digest[MD5_DIGEST_LENGTH];
  char cbuf[32];

  std::string hash, hash_mask, seq, seq_mask;
  std::string tile_path_step;

  m5_s.clear();

  for (allele=0; allele<2; allele++) {

    MD5_Init(&md5_ctx);

    tilestep = 0;
    while (tilestep < band_info.band[allele].size()) {

      if (tilestep >= sglf_path.ext_tileid.size()) { return -1; }
      if (tilestep >= sglf_path.seq.size()) { return -1; }

      span_len=1;
      while ( ((tilestep + span_len) < band_info.band[allele].size()) &&
              (band_info.band[allele][tilestep+span_len]==-1) ) {
        span_len++;
      }

      tilevar = band_info.band[allele][tilestep];

      if (tilevar >= sglf_path.ext_tileid[tilestep].size()) { return -2; }
      if (tilevar >= sglf_path.seq[tilestep].size()) { return -2; }

      seq = sglf_path.seq[tilestep][tilevar];

      noc_count=0;
      for (i=0; i<band_info.noc[allele][tilestep].size(); i+=2) {

        noc_start = band_info.noc[allele][tilestep][i];
        noc_len = band_info.noc[allele][tilestep][i+1];

        for (pos=noc_start; pos<(noc_start + noc_len); pos++) {
          seq[pos] = 'n';
        }

        noc_count += noc_len;

      }

      if (tilestep==0) {
        MD5_Update(&md5_ctx, (const void *)(seq.c_str()), (unsigned long)seq.size());
      } else {

        if (seq.size()>24) {
          MD5_Update(&md5_ctx, (const void *)(seq.c_str()+24), (unsigned long)(seq.size()-24));
        }

      }

      tilestep+=span_len;

    }

    MD5_Final(digest, &md5_ctx);

    if (allele>0) { m5_s += " "; }
    for (i=0; i<MD5_DIGEST_LENGTH; i++) {
      sprintf(cbuf, "%02x", (unsigned int)digest[i]);
      m5_s += cbuf;
    }

  }

  return 0;
}



int band_convert(band_info_t &band_info, sglf_path_t &sglf_path) {
  int i, j, k;
  int n, m, a;

  int allele=0;
  int tilestep=0, tilevar=0, span_len=0;
  int noc_count=0;
  int noc_start, noc_len, pos;
  int fold_w = 50;

  char *chp;

  std::string hash, hash_mask, seq, seq_mask;
  std::string tile_path_step;

  for (allele=0; allele<2; allele++) {

    tilestep = 0;
    while (tilestep < band_info.band[allele].size()) {

      if (tilestep >= sglf_path.ext_tileid.size()) { return -1; }
      if (tilestep >= sglf_path.seq.size()) { return -1; }

      span_len=1;
      while ( ((tilestep + span_len) < band_info.band[allele].size()) &&
              (band_info.band[allele][tilestep+span_len]==-1) ) {
        span_len++;
      }

      tilevar = band_info.band[allele][tilestep];

      if (tilevar >= sglf_path.ext_tileid[tilestep].size()) { return -2; }
      if (tilevar >= sglf_path.seq[tilestep].size()) { return -2; }

      seq = sglf_path.seq[tilestep][tilevar];
      seq_mask = sglf_path.seq[tilestep][tilevar];

      noc_count=0;
      for (i=0; i<band_info.noc[allele][tilestep].size(); i+=2) {

        noc_start = band_info.noc[allele][tilestep][i];
        noc_len = band_info.noc[allele][tilestep][i+1];

        for (pos=noc_start; pos<(noc_start + noc_len); pos++) {

          seq[pos] = 'n';
          if ((pos<24) || (pos > (seq.size()-24))) {
            seq_mask[pos] -= 'a' - 'A';
          }
          else { seq_mask[pos] = 'n'; }
        }

        noc_count += noc_len;

      }
      md5str(hash, seq);
      md5str(hash_mask, seq_mask);


      tile_path_step.clear();
      n=0;
      for (chp = (char *)sglf_path.ext_tileid[tilestep][tilevar].c_str();
           *chp; chp++) {
        if (*chp == '.') { n++; }
        if (n==3) { break; }
        tile_path_step += *chp;
      }

      chp = strchr((char *)sglf_path.ext_tileid[tilestep][tilevar].c_str(), '+');
      n = (int)sglf_path.ext_tileid[tilestep][tilevar].size();
      if (chp) {
        n = (int)(chp - sglf_path.ext_tileid[tilestep][tilevar].c_str());
      }

      // Print FastJ header
      printf(">{");
      //printf("\"%s\":\"%s\",", "tileID", sglf_path.ext_tileid[tilestep][tilevar].c_str());
      //printf("\"%s\":\"", "tileID"); print_substr(sglf_path.ext_tileid[tilestep][tilevar], 0, n); printf("\",");
      printf("\"%s\":\"%s.%03x\",", "tileID", tile_path_step.c_str(), allele);
      printf("\"%s\":\"%s\",", "md5sum", hash.c_str());
      printf("\"%s\":\"%s\",", "tagmask_md5sum", hash_mask.c_str());
      printf("\"%s\":%s,", "locus", "[ ]");
      printf("\"%s\":%i,", "n", (int)sglf_path.seq[tilestep][tilevar].size());
      printf("\"%s\":%i,", "seedTileLength", span_len);
      printf("\"%s\":%s,", "startTile", (tilestep==0) ? "true" : "false" );
      printf("\"%s\":%s,", "endTile", ((tilestep+span_len)==band_info.band[allele].size()) ? "true" : "false" );
      printf("\"%s\":\"", "startSeq"); print_substr(seq, 0, 24); printf("\",");
      printf("\"%s\":\"", "endSeq"); print_substr(seq, seq.size()-24, 24); printf("\",");
      printf("\"%s\":\"", "startTag"); print_substr(sglf_path.seq[tilestep][tilevar], 0, 24); printf("\",");
      printf("\"%s\":\"", "endTag"); print_substr(sglf_path.seq[tilestep][tilevar], sglf_path.seq[tilestep][tilevar].size()-24, 24); printf("\",");
      printf("\"%s\":%i,", "nocallCount", noc_count);
      printf("\"%s\":%s", "notes", "[ ]");
      printf("}\n");

      print_fold(seq, fold_w);
      printf("\n");

      tilestep+=span_len;

    }

  }

}

void print_sglf2bit_tilepath(sglf2bit_tilepath_t &sglf2bit_tilepath) {
  size_t n, m, tilestep, tilevar;
  std::string seq;

  n = sglf2bit_tilepath.ext_tileid.size();
  if ((n!=sglf2bit_tilepath.hash.size()) ||
      (n!=sglf2bit_tilepath.seq2bit.size())) {
    fprintf(stderr, "ERROR: size mismatch in sglf2bit_tilepath %i != (%i, %i)\n",
        (int)n,
        (int)sglf2bit_tilepath.hash.size(),
        (int)sglf2bit_tilepath.seq2bit.size());
    return;
  }

  for (tilestep=0; tilestep<n; tilestep++) {
    m = sglf2bit_tilepath.ext_tileid[tilestep].size();

    if ((m!=sglf2bit_tilepath.hash[tilestep].size()) ||
        (m!=sglf2bit_tilepath.seq2bit[tilestep].size())) {
      fprintf(stderr, "ERROR: size mismatch on tilestep %i in sglf2bit_tilepath %i != (%i, %i)\n",
          (int)tilestep,
          (int)m,
          (int)sglf2bit_tilepath.hash[tilestep].size(),
          (int)sglf2bit_tilepath.seq2bit[tilestep].size());
      return;
    }

    for (tilevar=0; tilevar<m; tilevar++) {
      seq.clear();
      sglf2bit_tilepath.seq2bit[tilestep][tilevar].twoBitToDnaSeq(seq);

      printf("%s,%s,%s\n",
          sglf2bit_tilepath.ext_tileid[tilestep][tilevar].c_str(),
          sglf2bit_tilepath.hash[tilestep][tilevar].c_str(),
          seq.c_str());

    }

  }
}

int batch_fastj_to_band(FILE *ifp_stream, FILE *sglf_fp) {
  sglf2bit_tilepath_t sglf2bit_tilepath;
  std::vector< fj_tile_t > fj_tile;
  band_info_t band_info;
  int ch, ret;
  gzFile gz_fp;
  FILE *ifp;
  std::string fn;

  fn.clear();

  read_sglf2bit_tilepath(sglf_fp, sglf2bit_tilepath);

  while (!feof(ifp_stream)) {
    ch = fgetc(ifp_stream);
    if ((ch==EOF) || (ch=='\n')) {
      if (fn.size()==0) { continue; }

      band_info.band[0].clear();
      band_info.band[1].clear();
      band_info.noc[0].clear();
      band_info.noc[1].clear();

      gz_fp = gzopen(fn.c_str(), "r");
      if (gz_fp == Z_NULL) { return -1; }
      fn.clear();

      read_tiles_gz(gz_fp, fj_tile);

      gzclose(gz_fp);

      ret = create_band_info_sglf2bit(band_info, fj_tile, sglf2bit_tilepath);
      if (ret<0) { return ret; }

      band_print(band_info);

      cleanup_tiles(fj_tile);

      fj_tile.clear();

      continue;
    }

    fn += (char)ch;
  }

  return 0;
}

int fastj_check (std::vector< fj_tile_t > &fj_tile) {
  int i, j, k;
  int count = 0;
  int seqSize;
  int tagSize;
  std::string seqHash;

  for (i = 0; i<fj_tile.size(); i++){
    if ( (int)fj_tile[i].seq.size() != fj_tile[i].n ) { return -1; }
    count = 0;
    for (j = 0; j<fj_tile[i].seq.size(); j++) {
      if ( (fj_tile[i].seq[j] == 'n') || (fj_tile[i].seq[j] == 'N') ) { count++; }
    }
    if ( count != fj_tile[i].nocall ) { return -2; }

    if ( fj_tile[i].startTag.size() > fj_tile[i].seq.size() ) { return -3; }
    for (j = 0; j<fj_tile[i].startTag.size(); j++) {
      if ( !((fj_tile[i].startTag[j] == fj_tile[i].seq[j]) ||
            (fj_tile[i].seq[j] == 'n') ||
            (fj_tile[i].seq[j] == 'N') )) { return -4; } // -3 already used
    }

    // Get sequence size and endTag size for convenience, then work backward
    // through string comparing them. endTag should be 24 or 0, but we can't
    // assume it. No calls are acceptable.
    //
    seqSize = fj_tile[i].seq.size();
    tagSize = fj_tile[i].endTag.size();
    if ( tagSize > seqSize ) { return -5; }
    for (j = 0; j<tagSize; j++) {
      if ( ! ((fj_tile[i].endTag[j] == fj_tile[i].seq[seqSize - tagSize + j])
           || (fj_tile[i].seq[seqSize - tagSize + j] == 'n')
           || (fj_tile[i].seq[seqSize - tagSize + j] == 'N') )) { return -6; }
    }

    md5str( seqHash, fj_tile[i].seq );
    if ( seqHash != fj_tile[i].seqHash ) { return -7; }

    if ( fj_tile[i].span <= 0 ) { return -8; }


  }
  return 0;
}

int fastj_check_tileid (std::vector< fj_tile_t > &fj_tile) {
  int i, j, k;
  int count = 0;
  int expectedTileStep[2] = {0,0};
  int expectedEndTileStep = -1;
  int endCount = 0;
  uint16_t ts, tv;

  // If the size is zero, there will be a out-of-bounds error
  //
  if ( fj_tile.size() == 0 ) { return -9; }

  ts = tileid_part(fj_tile[fj_tile.size()-1].tileid, 1);
  expectedEndTileStep = (int)ts + fj_tile[fj_tile.size()-1].span;

  for (i = 0; i<fj_tile.size(); i++){
    tv = tileid_part(fj_tile[i].tileid, 0);
    ts = tileid_part(fj_tile[i].tileid, 1);
    if (! ((tv == 0) || (tv == 1)) ) { return -10; }
    if ( ts != (uint16_t)expectedTileStep[tv] ) { return -11; }
    if ( (ts == 0) && (fj_tile[i].startTile != 1) ) { return -12; }
    expectedTileStep[tv] += fj_tile[i].span;
    if ( expectedTileStep[tv] == expectedEndTileStep ) {
      endCount++;
      if ( fj_tile[i].endTile != 1 ) { return -13; }
    }
  }
  if ( endCount != 2 ) { return -14; }

  return 0;
}

int main(int argc, char **argv) {
  int i, ret;
  int ch, opt, option_index;
  int check_tileid_option = 0;

  std::string ifn = "-", sglf_fn;

  int ifn_stream=0;
  std::vector< std::string > ifns;

  std::vector< fj_tile_t > fj_tile;
  sglf_path_t sglf_path;
  band_info_t band_info;
  std::vector< band_info_t > band_info_v;
  int show_help_flag = 1;
  int unsorted_flag = 0;

  int fold_width = 50;
  int tilepath=-1;

  std::string seq;
  std::string m5_s;

  uint64_t u64;
  uint16_t variant_id;

  FILE *ifp=stdin, *sglf_fp=NULL;

  FJT_ACTION action = FJT_NOOP;

  while ((opt=getopt_long(argc, argv, "vVhc:CL:i:p:BbHUI:Tt", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "invalid option, exiting\n");
      exit(-1);
      break;

    case 'T':
      show_help_flag=0;
      action = FJT_TEST;
      break;

    case 't':
      show_help_flag=0;
      check_tileid_option = 1;
      break;

    case 'C':
      show_help_flag=0;
      action = FJT_CSV;
      break;

    case 'U':
      unsorted_flag = 1;
      break;

    case 'c':
      show_help_flag=0;
      action = FJT_CONCAT;
      variant_id = (uint16_t)atoi(optarg);
      break;

    case 'L':
      sglf_fn = optarg;
      show_help_flag=0;
      break;

    case 'i':
      show_help_flag=0;
      ifn = optarg;
      break;

    case 'I':
      show_help_flag=0;
      ifn_stream = 1;
      ifn = optarg;
      break;

    case 'p':
      show_help_flag=0;
      tilepath = atoi(optarg);
      break;

    case 'B':
      show_help_flag=0;
      action = FJT_BAND;
      break;

    case 'b':
      show_help_flag=0;
      action = FJT_BAND_CONVERT;
      break;

    case 'H':
      show_help_flag=0;
      action = FJT_BAND_BATCH_HASH;
      break;

    case 'v':
      show_help_flag=0;
      verbose_flag = 1;
      break;
    case 'V':
      show_help_flag=0;
      show_version();
      exit(0);
      break;

    case 'h':
    default:
      show_help();
      exit(0);
      break;
  }

  if (show_help_flag) { show_help(); exit(0); }

  if (argc>optind) { ifn = argv[optind]; }
  if ((ifn.size()>0) && (ifn!="-")) {
    if ((ifp = fopen(ifn.c_str(), "r")) == NULL) {
      perror(ifn.c_str());
      exit(-1);
    }
  }

  if (action==FJT_CSV) {
    if (unsorted_flag==0) {
      read_tiles(ifp, fj_tile);
    }
    else {
    }
  }
  else if ( (action != FJT_BAND_CONVERT) &&
            (action != FJT_BAND_BATCH_HASH) ) {

    if ((action == FJT_BAND) &&
        (ifn_stream==1)) {
      // don't read in FastJ tiles.
      // The stream represents a list of files
      // that for batch processing with will be
      // read below.
      //
    }
    else {
      read_tiles(ifp, fj_tile);
    }

  }

  else if (action == FJT_BAND_CONVERT) {
    read_band(ifp, band_info);
  }

  else if (action == FJT_BAND_BATCH_HASH) {
    read_bands(ifp, band_info_v);
    //print_bands(band_info_v);
  }


  if (action ==  FJT_CSV) {

    if (unsorted_flag) {
      ret = read_tiles_and_print_csv(ifp);

      if (ret<0) {
        printf("ERROR: %i\n", ret);
      }

    }
    else {

      std::string m5;

      for (i=0; i<fj_tile.size(); i++) {

        print_tileid(fj_tile[i].tileid);
        printf("+%x", fj_tile[i].span);

        md5str(m5, fj_tile[i].seq);
        printf(",%s", m5.c_str());
        printf(",%s\n", fj_tile[i].seq.c_str());

      }
    }

  }

  else if (action == FJT_CONCAT) {
    concatenate_tiles(fj_tile, variant_id, seq);

    for (i=0; i<seq.size(); i++) {
      if ((i>0) && ((i%fold_width)==0)) { printf("\n"); }
      printf("%c", seq[i]);
    }
    printf("\n");
  }

  else if (action == FJT_TEST) {
    ret = fastj_check(fj_tile);
    if ( ret >= 0 ) {
      if ( verbose_flag ) {
        printf("fastj check: OK\n");
      }
    }
    else {
      fprintf(stderr, "fastj check error: %d\n", ret);
      exit(-1);
    }

    if (check_tileid_option) {
      ret = fastj_check_tileid(fj_tile);
      if ( ret >= 0 ) {
        if ( verbose_flag ) {
          printf("fastj tileid check: OK\n");
        }
      }
      else {
        fprintf(stderr, "fastj tileid check error: %d\n", ret);
        exit(-1);
      }
    }
  }

  else if (action == FJT_BAND) {

    if (sglf_fn.size()==0) {
      fprintf(stderr, "must provide SGLF file, exiting\n");
      exit(-1);
    }

    if ((sglf_fn == "-") && (ifp == stdin)) {
      fprintf(stderr, "SGLF stream must be different from FastJ input stream, exiting\n");
      exit(-2);
    }

    if (sglf_fn=="-") { sglf_fp = stdin; }
    else { sglf_fp = fopen(sglf_fn.c_str(), "r"); }
    if (!sglf_fp) {
      perror(sglf_fn.c_str());
      exit(-3);
    }

    if (ifn_stream) {

      ret = batch_fastj_to_band(ifp, sglf_fp);
      if (ret<0) {
        fprintf(stderr, "Error, invalid return code when converting batch band format: %i\n", ret);
        exit(-3);
      }

    }
    else {

      read_sglf_path(sglf_fp, sglf_path);

      ret = create_band_info(band_info, fj_tile, sglf_path);
      if (ret<0) {
        fprintf(stderr, "Error, invalid return code when covnerting to band format: %i\n", ret);
        exit(-3);
      }

      band_print(band_info);

    }

    if (sglf_fp!=stdin) { fclose(sglf_fp); }

  }

  else if (action == FJT_BAND_CONVERT) {

    if (sglf_fn.size()==0) {
      fprintf(stderr, "must provide SGLF file, exiting\n");
      exit(-1);
    }

    if ((sglf_fn == "-") && (ifp == stdin)) {
      fprintf(stderr, "SGLF stream must be different from FastJ input stream, exiting\n");
      exit(-2);
    }

    if (sglf_fn=="-") { sglf_fp = stdin; }
    else { sglf_fp = fopen(sglf_fn.c_str(), "r"); }
    if (!sglf_fp) {
      perror(sglf_fn.c_str());
      exit(-3);
    }

    read_sglf_path(sglf_fp, sglf_path);

    ret = band_convert(band_info, sglf_path);

    if (sglf_fp!=stdin) { fclose(sglf_fp); }

  }

  else if (action == FJT_BAND_BATCH_HASH) {

    if (sglf_fn.size()==0) {
      fprintf(stderr, "must provide SGLF file, exiting\n");
      exit(-1);
    }

    if ((sglf_fn == "-") && (ifp == stdin)) {
      fprintf(stderr, "SGLF stream must be different from FastJ input stream, exiting\n");
      exit(-2);
    }

    if (sglf_fn=="-") { sglf_fp = stdin; }
    else { sglf_fp = fopen(sglf_fn.c_str(), "r"); }
    if (!sglf_fp) {
      perror(sglf_fn.c_str());
      exit(-3);
    }

    read_sglf_path(sglf_fp, sglf_path);

    for (i=0; i<band_info_v.size(); i++) {

      ret = band_hash(m5_s, band_info_v[i], sglf_path);
      printf("%s\n", m5_s.c_str());
    }

    if (sglf_fp!=stdin) { fclose(sglf_fp); }


  }

  else if (action == FJT_FILTER) {
    fprintf(stderr, "not implemented...\n");
    exit(-1);
  }

  cleanup_tiles(fj_tile);
  if (ifp!=stdin) { fclose(ifp); }

  exit(0);
}
