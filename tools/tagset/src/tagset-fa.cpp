#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <getopt.h>

#include "bgzf.h"

#include <vector>
#include <string>

#define TF_VERSION "0.1.0"

inline static int _lc(int c) {
  if ((c>='A') && (c<='Z')) { return c - 'A' + 'a'; }
  return c;
}

typedef struct fai_type {
  std::vector< std::string > name;
  std::vector< long long int > beg, size, width, width_r;
} fai_t;

int split_tok_ch(std::vector< std::string > &tok_v, std::string &line, char field) {
  int i, j, k, n;
  std::string buf;

  n = (int)line.size();
  tok_v.clear();

  for (i=0; i<n; i++) {
    if (line[i]==field) {
      if (buf.size()>0) { tok_v.push_back(buf); }
      buf.clear();
      continue;
    }
    buf += line[i];
  }
  if (buf.size()>0) { tok_v.push_back(buf); }

  return 0;
}

long long int atoll(const char *s) {
  return strtoll(s, NULL, 10);
}

int read_start_pos(FILE *fp, std::vector< long long int > &start_pos) {
  int ch;
  std::string line;

  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch==EOF) || (ch=='\n') || (ch=='\r')) {
      if (line.size()==0) { continue; }
      start_pos.push_back(atoll(line.c_str()));
      line.clear();
      continue;
    }
    line += (char)ch;
  }

}

int print_tags(FILE *ofp, std::string &fa_fn, std::string &chrom, std::vector< long long int > &start_pos, int tag_len) {
  int i, k, ch, chrom_idx=-1, r;
  FILE *fp;
  BGZF *bgzfp;
  std::string idx_fn, line, buf;
  std::vector< std::string > tok_v;
  fai_t fai;
  long long int idx_offset, idx_size, idx_width, idx_width_r;
  long long int seq_q, seq_r, byte_offset;

  idx_fn =  fa_fn;
  idx_fn += ".fai";

  if (!(fp = fopen(idx_fn.c_str(), "r"))) { return -1; }
  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch==EOF) || (ch=='\n') || (ch=='\r')) {
      if (line.size()==0) { continue; }
      if (line[0]=='#') { continue; }

      split_tok_ch(tok_v, line, '\t');
      line.clear();
      if (tok_v.size()!=5) { continue; }

      fai.name.push_back(tok_v[0]);
      fai.beg.push_back(atoll(tok_v[2].c_str()));
      fai.size.push_back(atoll(tok_v[1].c_str()));
      fai.width.push_back(atoll(tok_v[3].c_str()));
      fai.width_r.push_back(atoll(tok_v[4].c_str()));

      if (tok_v[0] == chrom) {
        chrom_idx = (int)(fai.name.size()) - 1;
        idx_offset = fai.beg[chrom_idx];
        idx_size = fai.size[chrom_idx];
        idx_width = fai.width[chrom_idx];
        idx_width_r = fai.width_r[chrom_idx];
      }

      continue;
    }
    line += (char)ch;
  }
  fclose(fp);

  if (chrom_idx<0) { return -1; }

  if (!(bgzfp = bgzf_open(fa_fn.c_str(), "r"))) { return -2; }
  r = bgzf_index_load(bgzfp, fa_fn.c_str(), ".gzi");
  if (r<0) { return r; }

  for (i=0; i<start_pos.size(); i++) {

    seq_q = start_pos[i] / idx_width;
    seq_r = start_pos[i] % idx_width;

    byte_offset = idx_offset + (seq_q * idx_width_r) + seq_r;
    r = bgzf_useek(bgzfp, byte_offset, SEEK_SET);
    if (r<0) { return r; }

    buf.clear();
    //while ((byte_offset < end_offset) && (buf.size()<tag_len)) {
    while (buf.size()<tag_len) {
      byte_offset++;
      ch = bgzf_getc(bgzfp);
      if ((ch=='\n') || (ch=='\r')) { continue; }
      if (ch==EOF) { return -1; }
      buf += _lc(ch);
    }
    fprintf(ofp, "%s\n", buf.c_str());

  }

  bgzf_close(bgzfp);


}

static struct option long_options[] = {
    {"help",              no_argument,        NULL, 'h'},
    {"version",           no_argument,        NULL, 'v'},
    {"verbose",           no_argument,        NULL, 'V'},
    {"ref",               required_argument,  NULL, 'R'},
    {"chrom",             required_argument,  NULL, 'c'},
    {"input",             required_argument,  NULL, 'i'},
    {"output",            required_argument,  NULL, 'o'},
    {0,0,0,0}
};


void show_version() {
  printf("Version: %s\n", TF_VERSION);
}

void show_help() {
  printf("create tag sets from input start positions\n");
  printf("Version: %s\n", TF_VERSION);
  printf("usage:\n  tagsetFa [-h] [-v] [-V] [-i input_start_fn] [-o output_fa] [-R ref_fa]\n");
  printf("\n");
  printf("  [-i ifn]    input file of 0 reference start positions (default stdin)\n");
  printf("  [-o ofn]    output file of tag sequences (default stdout)\n");
  printf("  [-R ref_fa] reference FASTA file\n");
  printf("  [-c chrom]  chromosome name in FASTA file\n");
  printf("  [-V]        verbose\n");
  printf("  [-v]        version\n");
  printf("  [-h]        help (this screen)\n");
}

int main(int argc, char **argv) {
  int i, j, k, n, ch;
  std::vector< long long int > start_pos;
  std::string line;
  std::string ref_fa, ifn, ofn;
  FILE *ifp=NULL, *ofp=NULL;
  std::string name;
  int verbose_flag = 0;
  int opt, option_index=0;

  ifn = "-";
  ofn = "-";

  while ((opt = getopt_long(argc, argv, "hvVR:c:i:o:", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "sanity error, invalid option to parse, exiting\n");
      exit(-1);
      break;
    case 'R': ref_fa = optarg; break;
    case 'i': ifn = optarg; break;
    case 'o': ofn = optarg; break;
    case 'c': name = optarg; break;
    case 'V': verbose_flag = 1; break;
    case 'v': show_version(); exit(0); break;
    case 'h':
    default:
      show_help();
      exit(0);
      break;
  }

  if (argc>optind) {
    if ((argc-optind)==1) { ifn = argv[optind]; }
    else { fprintf(stderr, "invalid option"); show_help(); exit(-1); }
  }

  if (ref_fa.size()==0) { fprintf(stderr, "provide reference FASTA file\n"); show_help(); exit(-1); }
  if (ifn=="-") { ifp = stdin; }
  else {
    if (!(ifp = fopen(ifn.c_str(), "r"))) {
      perror(ifn.c_str()); exit(errno);
    }
  }

  read_start_pos(ifp, start_pos);
  if (ifp!=stdin) { fclose(ifp); }

  /*
  //DEBUG
  for (int ii=0; ii<start_pos.size(); ii++) {
    printf("[%i] %lli\n", ii, start_pos[ii]);
  }
  exit(1);
  //DEBUG
  */

  if (ofn=="-") { ofp = stdout; }
  else {
    if (!(ofp = fopen(ofn.c_str(), "w"))) {
      perror(ofn.c_str()); exit(errno);
    }
  }
  k = print_tags(ofp, ref_fa, name, start_pos, 24);

  if (ofp!=stdout) { fclose(ofp); }


}
