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

#include <vector>
#include <string>
#include <map>
#include <algorithm>

#include "cJSON.h"

#define FASTJ_TOOL_VERSION "0.1.0"

typedef struct fj_tile_type {
  cJSON *hdr;
  std::string seq;
  uint64_t tileid;
} fj_tile_t;

enum FJT_ACTION {
 FJT_NOOP = 0,
 FJT_CSV,
 FJT_CONCAT,
 FJT_FILTER
};

int verbose_flag = 0;

static struct option long_options[] = {
  {"help", no_argument, NULL, 'h'},
  {"verbose", no_argument, NULL, 'v'},
  {"version", no_argument, NULL, 'V'},
  {"csv", no_argument, NULL, 'C'},
  {"concatenate", required_argument, NULL, 'c'},
  {0,0,0,0}
};

void show_version() {
  printf("fjt version: %s\n", FASTJ_TOOL_VERSION);
}

void show_help() {
  show_version();
  printf("usage:\n  fjt [-c variant] [-C] [-v] [-V] [-h]\n");
  printf("\n");
  printf("  [-c variant]    Concatenate FastJ tiles into sequence.  `variant` is the variant ID to concatenate on.\n");
  printf("  [-C]            Output comma separated `tileID`, `hash` and `sequence` (CSV output).\n");
  printf("  [-v]            Verbose\n");
  printf("  [-V]            Version\n");
  printf("  [-h]            Help\n");
  printf("\n");
}

bool sortTileCmp(const fj_tile_t &lhs, const fj_tile_t &rhs) {
  return lhs.tileid < rhs.tileid;
}

enum fj_input_state { EXPECT_HDR, READ_HDR, READ_SEQ };

void print_tileid(uint64_t tileid) {
  uint64_t u64;
  int curpos=0;
  unsigned int byte_offset[] = { 6, 4, 2, 0 };

  for (curpos=0; curpos<4; curpos++) {
    u64 = tileid>>(8*byte_offset[curpos]);
    u64 &= 0xffff;
    if (curpos>0) { printf("."); }
    if (curpos!=1) {
      printf("%04x", (unsigned int)u64);
    } else {
      printf("%02x", (unsigned int)u64);
    }
  }

}

uint16_t tileid_part(uint64_t tileid, int part) {
  uint64_t u64;
  uint16_t u16;
  u64 = tileid>>(8*2*part);
  u64 &= 0xffff;
  u16 = (uint16_t)u64;
  return u16;
}

uint64_t parse_tileid(const char *tileid) {
  const char *chp;
  std::string s;
  unsigned long long int ull;
  uint64_t v=0, u64;
  int curpos=0;
  unsigned int byte_offset[] = { 6, 4, 2, 0 };

  for (chp=tileid; *chp; chp++) {
    if (*chp == '.') {
      ull = strtoull(s.c_str(), NULL, 16);
      u64 = (uint64_t)ull;
      v |= (u64 << (8*byte_offset[curpos]));
      curpos++;
      if (curpos>=4) { break; }
      s.clear();
      continue;
    }

    s+=*chp;
  }

  if (curpos<4) {
    ull = strtoull(s.c_str(), NULL, 16);
    u64 = (uint64_t)ull;
    v |= (u64 << (8*byte_offset[curpos]));
  }

  return v;
}

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
    } else {
      //printf("ERROR, element %i does not have tileID\n", i);
      return -4;
    }
  }

  std::sort( fj_tile.begin(), fj_tile.end(), sortTileCmp );

  /*
  printf("....\n");
  for (i=0; i<fj_tile.size(); i++) {
    cJSON *tid = cjson_obj(fj_tile[i].hdr, "tileID");
    if (cJSON_IsString(tid)) {
      printf("%s,%016llx,", tid->valuestring, (unsigned long long int)fj_tile[i].tileid);
    } else {
      printf("ERROR, element %i does not have tileID\n", i);
      return -4;
    }
    printf("%s", fj_tile[i].seq.c_str());
    printf("\n");
  }


  //DEBUG
  for (i=0; i<fj_tile.size(); i++) { cJSON_Delete(fj_tile[i].hdr); }
  */

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

int main(int argc, char **argv) {
  int i;
  int ch, opt, option_index;
  std::string ifn;
  std::vector< fj_tile_t > fj_tile;
  int show_help_flag = 1;

  int fold_width = 50;

  std::string seq;

  uint64_t u64;
  uint16_t variant_id;

  FILE *ifp=stdin;

  FJT_ACTION action = FJT_NOOP;

  while ((opt=getopt_long(argc, argv, "vVhc:C", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "invalid option, exiting\n");
      exit(-1);
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
    case 'C':
      show_help_flag=0;
      action = FJT_CSV;
      break;
    case 'c':
      show_help_flag=0;
      action = FJT_CONCAT;
      variant_id = (uint16_t)atoi(optarg);
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

  read_tiles(ifp, fj_tile);

  if (action ==  FJT_CSV) {
    std::string m5;

    for (i=0; i<fj_tile.size(); i++) {
      print_tileid(fj_tile[i].tileid);

      md5str(m5, fj_tile[i].seq);
      printf(",%s", m5.c_str());
      printf(",%s\n", fj_tile[i].seq.c_str());
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

  else if (action == FJT_FILTER) {
    printf("not implemented...\n");
  }

  cleanup_tiles(fj_tile);
  if (ifp!=stdin) { fclose(ifp); }
}
