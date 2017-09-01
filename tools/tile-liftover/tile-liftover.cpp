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

// Sample usage:
//
//   refstream chrM | ./tile-liftover -T <( refstream /data-sdd/data/l7g/tagset.fa/tagset.fa.gz 035e.00 ) -c chrM -p 862 -N hg19
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <errno.h>
#include <getopt.h>

#include <vector>
#include <string>
#include <map>

extern "C" {
#include "asm_ukk.h"
}

#define TILE_LIFTOVER_VERSION "0.2.0"

int verbose_flag = 0;

static struct option long_options[] = {
  {"help",            no_argument, NULL, 'h'},
  {"verbose",         no_argument, NULL, 'v'},
  {"version",         no_argument, NULL, 'V'},
  {"greedy-match",    no_argument, NULL, 'G'},
  {"optimize-match",  no_argument, NULL, 'O'},
  {"tagset",          required_argument, NULL, 'T'},
  {"tilepath",        required_argument, NULL, 'p'},
  {"endseq",          required_argument, NULL, 'E'},
  {"end-tile-length", required_argument, NULL, 'M'},
  {"ref-stream",      required_argument, NULL, 'R'},
  {"chrom",           required_argument, NULL, 'c'},
  {"ref-name",        required_argument, NULL, 'N'},
  {"start",           required_argument, NULL, 's'},
  {0,0,0,0}
};

void show_version() {
  printf("tile-liftover version: %s\n", TILE_LIFTOVER_VERSION);
}


void show_help() {
  show_version();
  printf("usage:\n");
  printf("  -T tagset       tagset stream\n");
  printf("  -p tilepath     tilepath\n");
  printf("  [-R ref]        reference stream (stdin default)\n");
  printf("  [-s start]      start position (0 reference, 0 default)\n");
  printf("  [-c chrom]      chromosome\n");
  printf("  [-G]            greedy match\n");
  printf("  [-O]            optimize alignment (default)\n");
  printf("  [-M M]          don't go past M bases from the last tag in the reference stream\n");
  printf("  [-E endseq]     End on endseq (takes precedence over M)\n");
  printf("  [-N refname]    reference name (default 'hg19')\n");
  printf("  [-v]            verbose\n");
  printf("  [-V]            version\n");
  printf("  [-h]            help\n");
  printf("\n");
}

// Our dynamic programming match function.
// mismatch will be set to something large
// so that mismatches are not allowed.
//
int g_mismatch = 100000;
int g_gap = 1;

extern "C" {
  int score_func(void *x, void *y, size_t sz) {
    int *ix, *iy;
    ix = (int *)x;
    iy = (int *)y;

    if ((!ix) || (!iy)) { return g_gap; }
    if ( (*ix) == (*iy) ) { return 0; }
    return g_mismatch;

  }
}

// Try to optimally match tags as they're found in the reference stream.
// Find an alignment of the contiguous tile steps to the found tile steps,
// not allowing mismatches.
//
// If `the end_tile_length` is greater than 0, only take `end_tile_length` (tag
// inclusive) of the last tile from the reference stream.
// If `end_seq` is non-empty, stop if the end-sequence is found.
// `end_seq` takes precedence over `end_tile_length`
//
// Note that this does an dymaic programming alignment of the two sequences.  For
// normal use cases this should be fine but if the two alignments are too far
// off (something went wrong, say) then memory usage and run-time could become
// prohibitive.
//

int match_tag(FILE *ref_fp,
              std::string &ref_name, std::string &chrom_str, int start_pos, int tilepath,
              std::string &tagset, std::map< std::string, int > &tag_pos_map, std::string &end_seq, int end_tile_length) {
  int n_tag=0, count, i, j, k, tag_idx;
  int ch;
  std::string ref, tag;
  std::map< std::string, int >::iterator tpm_it;

  std::vector< int > orig_tile_step, match_tile_step;
  int *X, *Y;
  size_t X_len, Y_len;
  int sc, gap_int = -1;

  int orig_n_tag=-1;
  const char *chp;

  int64_t end_pos_non_inc=-1;

  std::vector< int > final_tag;
  char fmt_str[] = "%04x\t%10i\n";

  if ((tagset.size()%24)!=0) {
    return -1;
  }
  orig_n_tag = tagset.size()/24;
  if (orig_n_tag==0) {
    printf(fmt_str, 0, (int)(start_pos + end_tile_length));
    return -1;
  }

  while (!feof(ref_fp)) {
    ch = fgetc(ref_fp);
    if (ch==EOF) { continue; }
    if ((ch=='\r') || (ch=='\n') || (ch==' ')) { continue; }
    ref += (char)ch;
  }

  for (tpm_it = tag_pos_map.begin(); tpm_it != tag_pos_map.end(); ++tpm_it) {
    if (n_tag < tpm_it->second) {
      n_tag = tpm_it->second;
    }
  }

  n_tag++;
  for (i=0; i<n_tag; i++) {
    orig_tile_step.push_back(i);
  }

  end_pos_non_inc = -1;
  count=0;
  for (i=23; i<ref.size(); i++) {
    tag.clear();
    for (k=0; k<24; k++) { tag += ref[i-23+k]; }
    if ( tag_pos_map.find(tag) != tag_pos_map.end() ) {

      //DEBUG
      //printf("match! %i %s\n", count, tag.c_str());

      match_tile_step.push_back( tag_pos_map[tag] );
      count++;


      if (end_tile_length>0) {

        /*
        // Check against last tag
        //
        if (strncmp(tag.c_str(), tagset.c_str() + tagset.size() - 24, 24)==0) {
          end_pos_non_inc = (int64_t)(i - 23 + end_tile_length);
        }

        // Force the end position to be end_tile_length past
        // the last observed tile (rather than the last tile in
        // the tilepath)
        //
        end_pos_non_inc = (int64_t)(i-23+end_tile_length);
        */
      }

    }
  }

  if (end_seq.size()>0) {
    chp = strstr(ref.c_str(), end_seq.c_str());
    if (chp) {
      end_pos_non_inc = (int64_t)(chp - &(ref[0]));
      end_pos_non_inc += (int64_t)end_seq.size();
    }
    else {
      //printf("end seq not found\n");
    }
  }

  // don't allow any mismatches
  //
  g_mismatch = 10*n_tag;

  //DEBUG
  printf("ALIGNING:\n");
  for (i=0; i<orig_tile_step.size(); i++) {
    printf(" %i", orig_tile_step[i]);
  }
  printf("\n");
  for (i=0; i<match_tile_step.size(); i++) {
    printf(" %i", match_tile_step[i]);
  }
  printf("\n");
  fflush(stdout);
  //DEBUG

  sc =
    avm_ukk_align3((void **)(&X), &X_len,
                   (void **)(&Y), &Y_len,
                   (void *)(&(orig_tile_step[0])), orig_tile_step.size(),
                   (void *)(&(match_tile_step[0])), match_tile_step.size(),
                   score_func,
                   (void *)(&gap_int),
                   sizeof(int));

  for (i=0; i<(int)X_len; i++) {
    if ((X[i] < 0) || (Y[i] < 0)) { continue; }
    if (X[i] != Y[i]) {
      printf("SANITY ERROR pos: X[%i] %i != Y[%i] %i\n", i, X[i], i, Y[i]);
      return -1;
      continue;
    }
    final_tag.push_back(X[i]);
  }

  //DEBUG
  //printf("X:"); for (i=0; i<X_len; i++) { printf(" %4i", X[i]); } printf("\n");
  //printf("Y:"); for (i=0; i<Y_len; i++) { printf(" %4i", Y[i]); } printf("\n");
  //printf("fin:"); for (i=0; i<final_tag.size(); i++) { printf(" %i", final_tag[i]); } printf("\n");
  //DEBUG

  if (X) { free(X); }
  if (Y) { free(Y); }

  if (final_tag.size()==0) { return -1; }

  tag_idx = 0;
  tag.clear();
  for (j=0; j<24; j++) { tag += tagset[ final_tag[tag_idx]*24 + j ]; }

  printf(">%s:%s:%04x\n", ref_name.c_str(), chrom_str.c_str(), tilepath);

  if (end_pos_non_inc<0) { end_pos_non_inc = (int)ref.size(); }
  //for (i=23; i<ref.size(); i++) {
  for (i=0; i<(ref.size()-24); i++) {

    if (strncmp(tag.c_str(), ref.c_str() + i, 24)==0) {
      printf(fmt_str, final_tag[tag_idx], i+24 + start_pos);
      tag_idx++;
      tag.clear();

      if (end_tile_length > 0) {
        end_pos_non_inc = (int64_t)(i+end_tile_length);
      }

      if (tag_idx >= final_tag.size()) { break; }
      for (j=0; j<24; j++) { tag += tagset[ final_tag[tag_idx]*24 + j ]; }

    }

  }

  if (tag_idx==0) {
    printf(fmt_str, 1, (int)(end_pos_non_inc + start_pos));
    return -1;
  }

  printf(fmt_str, final_tag[tag_idx-1]+1, (int)(end_pos_non_inc + start_pos));
}

//
// Simple tag liftover:
// print out first matching tag as found in `tag_pos_map`.
//
// Print a warning on skipped tags or duplicate tags.
//

int greedy_match(FILE *ref_fp,
                 std::string &ref_name, std::string &chrom_str, int start_pos,
                 int tilepath,
                 std::map< std::string, int > &tag_pos_map) {
  int i, pos, idx_pos;
  int cur_tag_id, ch;
  std::string tag, ref_str;
  char fmt_str[] = "%04x\t%10i\n";


  // Read in reference stream, recording the position where we find each of the tags.
  //

  pos = start_pos;

  printf(">%s:%s:%04x\n", ref_name.c_str(), chrom_str.c_str(), tilepath);

  cur_tag_id = -1;
  idx_pos = 0;
  std::string prev_tag;
  while (!feof(ref_fp)) {
    ch = fgetc(ref_fp);
    if (ch==EOF) { continue; }
    if ((ch=='\n') || (ch==' ') || (ch=='\r')) { continue; }

    ref_str += (char)ch;
    if (ref_str.size() >= 24) {
      tag.clear();
      for (i=0; i<24; i++) {
        tag += ref_str[ idx_pos - 24 + i ];
      }

      if ( tag_pos_map.find(tag) != tag_pos_map.end() ) {

        if (cur_tag_id >= tag_pos_map[tag]) {
          printf("WARNING: tag %s (tag#:%i) found before current end tag %s (tag#:%i), skipping\n",
              tag.c_str(), tag_pos_map[tag], prev_tag.c_str(), cur_tag_id);
        }
        else {
          if (verbose_flag) { printf("# %s %i %i\n", tag.c_str(), tag_pos_map[tag], pos); }
          printf(fmt_str, tag_pos_map[tag], pos);
          cur_tag_id = tag_pos_map[tag];

          prev_tag.clear();
          prev_tag = tag;
        }

      }
    }

    pos++;
    idx_pos++;

  }

  return 0;
}



int main(int argc, char **argv) {
  int i, j, k;
  int opt, ch, option_index;
  int pos, start_pos=0;
  int end_tilestep=-1;
  FILE *ref_fp, *ofp;
  FILE *tagset_fp=NULL;

  int greedy_match_flag = 0;
  int opt_match_flag = 1;

  std::string ref_fn = "-";
  std::string tagset_fn;
  std::string tagset;
  std::map< std::string, int > tag_pos_map;
  std::string tag;
  std::string ref_str;

  std::string ref_name;
  std::string chrom_str;
  int tilepath=-1;

  char fmt_str[] = "%04x\t%10i\n";

  std::string end_seq;
  int end_tile_length=-1;

  ref_fp = stdin;
  ofp = stdout;

  ref_name = "hg19";
  chrom_str = "unk";

  while ((opt = getopt_long(argc, argv, "T:s:c:p:R:N:vVhE:M:GO", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "sanity error, invalid option to parse, exiting\n");
      exit(-1);
      break;
    case 'T':
      tagset_fn = optarg; break;
    case 'p':
      tilepath = atoi(optarg); break;
    case 'c':
      chrom_str = optarg; break;
    case 'N':
      ref_name = optarg; break;
    case 'v':
      verbose_flag = 1; break;
    case 'V':
      show_version(); exit(0); break;
    case 's':
      start_pos = atoi(optarg); break;
    case 'E':
      end_seq = optarg; break;
    case 'M':
      end_tile_length = atoi(optarg); break;
    case 'R':
      ref_fn = optarg; break;
    case 'G':
      greedy_match_flag = 1; opt_match_flag = 0; break;
    case 'O':
      greedy_match_flag = 0; opt_match_flag = 1; break;
    case 'h':
    default:
      show_help();
      exit(0);
      break;
  }

  if (tagset_fn.size()==0) {
    fprintf(stderr, "provide tagset\n");
    show_help();
    exit(-1);
  }

  if (tagset_fn == "-") { tagset_fp = stdin; }
  else {
    tagset_fp = fopen(tagset_fn.c_str(), "r");
    if (tagset_fp==NULL) {
      perror(tagset_fn.c_str());
      exit(-1);
    }
  }

  if (ref_fn.size()==0) {
    fprintf(stderr, "invalid reference stream\n");
    show_help();
    exit(-1);
  }

  if (ref_fn=="-") { ref_fp = stdin; }
  else {
    ref_fp = fopen(ref_fn.c_str(), "r");
    if (ref_fp==NULL) {
      perror(ref_fn.c_str());
      exit(-1);
    }
  }

  if (tagset_fp == ref_fp) {
    fprintf(stderr, "tagset and reference streams are identical, exiting\n");
    show_help();
    exit(-2);
  }

  if (tilepath == -1) {
    fprintf(stderr, "provide tilepath\n");
    show_help();
    exit(-3);
  }

  // Read tagset
  //

  pos = 0;

  if (verbose_flag) { printf("# reading tagset..."); fflush(stdout); }

  while (!feof(tagset_fp)) {
    ch = fgetc(tagset_fp);
    if (ch==EOF) { continue; }
    if ((ch=='\n') || (ch==' ')) { continue; }
    tagset += (char)ch;

    if (tag.size()==24) {
      tag_pos_map[tag] = (pos-24)/24;
      tag.clear();
    }

    tag += (char)ch;
    pos++;
  }

  if (tag.size()==24) {
    tag_pos_map[tag] = (pos-24)/24;
    end_tilestep = (pos-24)/24;
    tag.clear();
  }

  if (verbose_flag) { printf("done\n"); fflush(stdout); }

  if ((tagset.size()%24)!=0) {
    printf("Incorrect tagset size (%i).  Expecting exact divisor of tag length (24)\n", (int)tagset.size());
    exit(-3);
  }

  if (opt_match_flag) {

    match_tag(ref_fp,
              ref_name, chrom_str, start_pos, tilepath,
              tagset, tag_pos_map, end_seq, end_tile_length);

  }
  else {

    greedy_match(ref_fp,
                 ref_name, chrom_str, start_pos, tilepath,
                 tag_pos_map);

  }

  if (ref_fp!=stdin) { fclose(ref_fp); }
  if (tagset_fp!=stdin) { fclose(tagset_fp); }
  if (ofp!=stdout) { fclose(ofp); }
  return 0;
}
