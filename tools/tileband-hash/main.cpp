#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>

#include <errno.h>
#include <openssl/md5.h>

#include <vector>
#include <string>

#include "sglf.hpp"
#include "tileband-hash.hpp"

#define TILEBAND_HASH_VERSION "0.1.0"

static struct option long_options[] = {
  {"help", no_argument, NULL, 'h'},
  {"verbose", no_argument, NULL, 'v'},
  {"version", no_argument, NULL, 'V'},

  {"n-dataset", required_argument, NULL, 'n'},
  {"sglf", required_argument, NULL, 'L'},
  {"tilepaths", required_argument, NULL, 'T'},

  {0,0,0,0}
};

void show_version() {
  printf("tileband-hash version: %s\n", TILEBAND_HASH_VERSION);
}

void show_help() {
  show_version();
  printf("usage:\n  tileband-hash [-n N] [-L sglf_stream] [-T tilepaths] [-v] [-V] [-h] bands\n");
  printf("\n");
  printf("  -L sglf_stream  SGLF stream\n");
  printf("  -T tilepaths    decimal list of tilepaths (e.g. '752+2', '752-753', '752,753')\n");
  printf("  [-n N]          number of datasets to convert (input band count must be a 4x this number, default 1)\n");
  printf("  [-v]            Verbose\n");
  printf("  [-V]            Version\n");
  printf("  [-h]            Help\n");
  printf("\n");
}

int split_tok_ch(std::vector< std::string > &tok_v, std::string &line, const char *fields) {
  int i, j, k, n, n_fields;
  std::string buf;
  char default_fields[] = "\t";
  int subsume_whitespace = 1;

  n = (int)line.size();
  tok_v.clear();

  if (!fields) {
    fields = default_fields;
  }
  n_fields = strlen(fields);

  for (i=0; i<n; i++) {
    for (j=0; j<n_fields; j++) {
      if (line[i]==fields[j]) {
        if (buf.size()>0) { tok_v.push_back(buf); }
        buf.clear();
        break;
      }
    }

    if (j<n_fields) {

      if (subsume_whitespace) {
        while ((i+1)<n) {
          for (j=0; j<n_fields; j++) {
            if (line[i+1] == fields[j]) { break; }
          }
          if (j==n_fields) { break; }
          i++;
        }

      }

      continue;
    }
    buf += line[i];
  }
  if (buf.size()>0) { tok_v.push_back(buf); }

  return 0;
}



int parse_tilepaths(std::vector<int> &tilepath_list, std::string &tilepaths) {
  int i, j, k;
  int r;

  long int s, n, e;

  std::vector< std::string > ele_list, val_list;

  r = split_tok_ch(ele_list, tilepaths, ",");
  if (r<0) { return r; }

  for (i=0; i<ele_list.size(); i++) {
    if (strchr(ele_list[i].c_str(), '+')) {
      r = split_tok_ch(val_list, ele_list[i], "+");
      if (r<0) { return r; }
      if (val_list.size()!=2) { return -2; }

      s = strtol(val_list[0].c_str(), NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -3; }

      n = strtol(val_list[1].c_str(), NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -3; }

      if (s<0) { return -4; }
      if (n<0) { return -5; }

      for (k=s; k<(s+n); k++) { tilepath_list.push_back(k); }

    }

    else if (strchr(ele_list[i].c_str(), '-')) {
      r = split_tok_ch(val_list, ele_list[i], "-");
      if (r<0) { return r; }
      if (val_list.size()!=2) { return -2; }

      s = strtol(val_list[0].c_str(), NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -3; }

      e = strtol(val_list[1].c_str(), NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -3; }

      if (s<0) { return -4; }
      if (e<s) { return -5; }

      for (k=s; k<e; k++) { tilepath_list.push_back(k); }
    }

    else {
      s = strtol(ele_list[i].c_str(), NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -3; }
      tilepath_list.push_back(s);
    }
  }

  return 0;
}

int main(int argc, char **argv) {
  int i, ii;
  std::string ifn = "-", sglf_fn = "", tilepaths = "";
  FILE *ifp=NULL, *sglf_fp=NULL;
  int show_help_flag=1;
  int verbose_flag=0;
  int ch, opt, option_index;

  int n_dataset=1, r=0;

  std::vector< band_info_t > band_datasets;
  band_info_t band_info;
  std::vector< band_info_t > band_info_v;
  std::vector< int > tilepath_list;

  std::vector< std::string > digest;

  sglf_t sglf;

  while ((opt=getopt_long(argc, argv, "hvVL:T:n:", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "invalid option, exiting\n");
      exit(-1);
      break;

    case 'L':
      show_help_flag=0;
      sglf_fn = optarg;
      break;
    case 'T':
      show_help_flag=0;
      tilepaths = optarg;
      break;
    case 'n':
      show_help_flag=0;
      n_dataset=atoi(optarg);
      break;

    case 'V':
      show_help_flag=0;
      show_version();
      exit(0);
      break;
    case 'v':
      show_help_flag=0;
      verbose_flag=1;
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

  if (tilepaths.size()==0) {
    printf("please provide tilepaths\n\n");
    show_help();
    exit(0);
  }

  if (sglf_fn.size()==0) {
    printf("please provide sglf stream\n\n");
    show_help();
    exit(0);
  }

  if ((sglf_fp = fopen(sglf_fn.c_str(), "r")) == NULL) {
    perror(sglf_fn.c_str());
    exit(-1);
  }

  r = parse_tilepaths(tilepath_list, tilepaths);
  if (r<0) {
    printf("ERROR: could not parse tilepath list\n\n");
    show_help();
    exit(-1);
  }

  r = sglf_read(sglf_fp, sglf);
  if (r<0) {
    printf("ERROR: SGLF read: got %i\n", r);
    exit(r);
  }
  fclose(sglf_fp);

  r = read_bands(ifp, band_datasets);
  if (r<0) {
    printf("ERROR: band datset read: got %i\n", r);
    exit(r);
  }
  if (ifp!=stdin) { fclose(ifp); }

  if ( (band_datasets.size() % tilepath_list.size()) != 0 ){
  //if ( (int)(band_datasets.size()) != (n_dataset*(int)(tilepath_list.size())) ) {
    printf("ERROR: band dataset size mismatch: number bands %i, tilepath_list size %i, does not divide evenly\n",
        (int)band_datasets.size(),
        (int)(tilepath_list.size()));
    exit(-1);
  }

  band_md5_hash(digest, band_datasets, sglf, tilepath_list);

  for (i=0; i<digest.size(); i++)  {
    printf("%s\n", digest[i].c_str());
  }

}
