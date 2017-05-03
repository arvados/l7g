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

#include <errno.h>
#include <getopt.h>

#include <vector>
#include <string>
#include <map>

#define TILE_LIFTOVER_VERSION "0.1.0"

int verbose_flag = 0;

static struct option long_options[] = {
  {"help", no_argument, NULL, 'h'},
  {"verbose", no_argument, NULL, 'v'},
  {"version", no_argument, NULL, 'V'},
  {"tagset", required_argument, NULL, 'T'},
  {"tilepath", required_argument, NULL, 'p'},
  {"ref-stream", required_argument, NULL, 'R'},
  {"chrom", required_argument, NULL, 'c'},
  {"ref-name", required_argument, NULL, 'N'},
  {"start", required_argument, NULL, 's'},
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
  printf("  [-N refname]    reference name (default 'hg19')\n");
  printf("  [-v]            verbose\n");
  printf("  [-V]            version\n");
  printf("  [-h]            help\n");
  printf("\n");
}

int main(int argc, char **argv) {
  int i, j, k;
  int opt, ch, option_index;
  int pos, start_pos=0;
  int end_tilestep=-1;
  FILE *ref_fp, *ofp;
  FILE *tagset_fp=NULL;

  std::string tagset_fn;
  std::string tagset;
  std::map< std::string, int > tag_pos_map;
  std::string tag;
  std::string ref_str;

  std::string ref_name;
  std::string chrom_str;
  int tilepath=-1;

  char fmt_str[] = "%04x\t%10i\n";

  int show_help_flag = 1;

  ref_fp = stdin;
  ofp = stdout;


  ref_name = "hg19";
  chrom_str = "unk";

  while ((opt = getopt_long(argc, argv, "T:s:c:p:R:N:vVh", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "sanity error, invalid option to parse, exiting\n");
      exit(-1);
      break;
    case 'T':
      show_help_flag=0;
      tagset_fn = optarg; break;
    case 'p':
      show_help_flag=0;
      tilepath = atoi(optarg); break;
    case 'c':
      show_help_flag=0;
      chrom_str = optarg; break;
    case 'N':
      show_help_flag=0;
      ref_name = optarg; break;
    case 'v':
      show_help_flag=0;
      verbose_flag = 1; break;
    case 'V':
      show_help_flag=0;
      show_version(); exit(0); break;
    case 's':
      show_help_flag=0;
      start_pos = atoi(optarg); break;
    case 'h':
    default:
      show_help();
      exit(0);
      break;
  }

  if (show_help_flag) {
    show_help();
    exit(0);
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

  //DEBUG
  /*
  printf("??:\n%s\n", tagset.c_str());
  for (i=0; i<tagset.size(); i+=24) {
    printf("%3i:\n", i/24);
    for (j=0; j<24; j++) { printf("%c", tagset[i + j]); }
    printf("\n");
  }
  */


  // Read in reference stream, recording the position where we find each of the tags.
  //

  pos = start_pos;

  printf(">%s:%s:%04x\n", ref_name.c_str(), chrom_str.c_str(), tilepath);

  int cur_tag_id = -1;
  int idx_pos = 0;
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

  printf(fmt_str, end_tilestep+1, pos);

  if (ref_fp!=stdin) { fclose(ref_fp); }
  if (tagset_fp!=stdin) { fclose(tagset_fp); }
  if (ofp!=stdout) { fclose(ofp); }
  return 0;
}
