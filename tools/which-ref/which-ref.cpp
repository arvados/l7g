#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>

#include "asm_ukk.h"

#include <string>
#include <vector>

#define WHICH_REF_VERSION "0.1.0"

typedef struct opt_type {
  int debug;
  int show_all_score;
  int verbose;
  int version;
} opt_t;

static struct option long_options[] = {
  {"help",                no_argument,        NULL, 'h'},
  {"version",             no_argument,        NULL, 'v'},
  {"verbose",             no_argument,        NULL, 'V'},
  {0,0,0,0}
};


void print_usage() {
  printf("\n");
  printf("usage:\n");
  printf("\n");
  printf("    which-ref [-h] [-v] [-V] [query-seq0] [query-seq1] [query-seq2] ... [query-seqN] [ref-seq]\n");
  printf("\n");
  printf("  [-h]        Print help (this screen)\n");
  printf("  [-v]        Print version\n");
  printf("  [-V]        Verbose flag\n");
  printf("\n");
}

void print_version() {
  printf("version: %s\n", WHICH_REF_VERSION);
}

void init_opt(opt_t *opt) {
  opt->debug = 0;
  opt->show_all_score = 0;
  opt->verbose = 0;
}

int main(int argc, char **argv) {
  FILE *fp;
  int i, k;
  std::vector< std::string > ref_fns, ref_seq;
  std::string ifn, seq, tseq;
  int min_score, min_idx;
  int loc_debug = 0;

  int ch, option_index;
  opt_t opt;

  std::vector< int > score;

  while ((ch = getopt_long(argc, argv, "hvVA", long_options, &option_index))!=-1) switch(ch) {
    case 0:
      fprintf(stderr, "sanity error, invalid optino to parse, exiting\n");
      exit(-1);
      break;
    case 'v':
      print_version();
      exit(0);
      break;
    case 'V':
      opt.verbose=1;
      break;
    case 'A':
      opt.show_all_score=1;
      break;
    default:
    case 'h':
      print_usage();
      exit(0);
      break;
  }

  if (argc>optind) {
    for (i=0; i<(argc-optind-1); i++) {
      ref_fns.push_back( argv[optind+i] );
    }
    ifn = argv[optind+i];
  }

  if (ref_fns.size()==0) {
    printf("Provide input sequences to compare\n");
    print_usage();
    exit(0);
  }

  /*
  if (argc<3) {
    print_usage();
    exit(-1);
  }

  for (i=1; i<(argc-1); i++) {
  }

  ifn = argv[argc-1];
  */

  //DEBUG
  //
  if (loc_debug) {
    for (i=0; i<ref_fns.size(); i++) {
      printf("%s\n", ref_fns[i].c_str());
    }
    printf("ifn: %s\n", ifn.c_str());
  }

  for (i=0; i<ref_fns.size(); i++) {
    fp = fopen(ref_fns[i].c_str(), "r");
    if (!fp) {
      perror(ref_fns[i].c_str());
      exit(-1);
    }

    tseq.clear();
    while (!feof(fp)) {
      ch = fgetc(fp);
      if (ch==EOF)  { continue; }
      if (ch=='\n') { continue; }
      tseq += (char)ch;
    }
    ref_seq.push_back(tseq);

    fclose(fp);
  }

  fp = fopen(ifn.c_str(), "r");
  if (!fp) {
    perror(ifn.c_str());
    exit(-1);
  }

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF)  { continue; }
    if (ch=='\n') { continue; }
    seq += (char)ch;
  }
  fclose(fp);
  if (loc_debug) {
    for (i=0; i<ref_fns.size(); i++) {
      printf("seq(%s): %s\n",
          ref_fns[i].c_str(),
          ref_seq[i].c_str());
    }
    printf("inp: %s\n", seq.c_str());
  }

  min_score = -1;
  min_idx = 0;
  for (i=0; i<ref_fns.size(); i++) {
    k = asm_ukk_score( (char *)ref_seq[i].c_str(), (char *)seq.c_str() );
    score.push_back(k);

    if ((i==0) || (k<min_score)) {
      min_score = k;
      min_idx = i;
    }
  }

  printf("min_score: %i\nmin_idx:%i\nname:%s\n",
      min_score,
      min_idx,
      ref_fns[min_idx].c_str());

}
