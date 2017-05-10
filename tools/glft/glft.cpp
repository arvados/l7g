#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <getopt.h>
#include <openssl/md5.h>

#include <string>

#define GLF_TOOL_VERSION "0.1.0"

int verbose_flag = 0;

#define GLFT_NOOP 0
#define GLFT_PATH_LIB_VER 1
#define GLFT_LIB_VER 2
#define GLFT_LIB_MANIFEST 3
#define GLFT_PATH_LIB_MANIFEST 4
#define GLFT_N 5

static struct option long_options[] = {
  {"help", no_argument, NULL, 'h'},
  {"verbose", no_argument, NULL, 'v'},
  {"version", no_argument, NULL, 'V'},
  {"use-reported-hash", no_argument, NULL, 'H'},
  {"tile-lib-version", no_argument, NULL, 'L'},
  {"tile-lib-path-version", no_argument, NULL, 'P'},
  {0,0,0,0}
};

void show_version() {
  printf("glft version: %s\n", GLF_TOOL_VERSION);
}

void show_help() {
  show_version();
  printf("usage:\n");
  printf("  [-L]        Tile Library Version\n");
  printf("  [-P]        Tile Library Path Version. From an SGLF file, get path library version.\n");
  printf("  [-H]        Use hash as reported in SGLF file (default recalculate hash by sequence)\n");
  printf("  [-v]        Verbose\n");
  printf("  [-V]        Version\n");
  printf("  [-h]        Help\n");
  printf("\n");
}


typedef struct glft_opt_type {
  int use_reported_hash;
} glft_opt_t;


#define SGLF_RS_TILEID 0
#define SGLF_RS_HASH 1
#define SGLF_RS_SEQ 2
#define SGLF_RS_N 3


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


int create_tile_path_manifest(std::string &manifest, FILE *fp, glft_opt_t &glft_opt) {
  int ch;
  int read_state=0, line_no=0;
  std::string tileid_s, hash_s, seq_s, ts;
  std::string line;

  manifest.clear();

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {

      line_no++;

      if (line.size()==0) { continue; }

      /*
      printf("%s\n%s\n%s\n-----\n\n",
          tileid_s.c_str(),
          hash_s.c_str(),
          seq_s.c_str());
          */

      if (glft_opt.use_reported_hash) {
        ts = hash_s;
      } else {
        md5str(ts, seq_s);
      }

      if (manifest.size()>0) { manifest += ' '; }
      manifest += tileid_s;
      manifest += ':';
      manifest += ts;

      tileid_s.clear();
      hash_s.clear();
      seq_s.clear();
      read_state = 0;

      line.clear();
      continue;
    }

    if (ch==',') {
      read_state = (read_state+1)%(SGLF_RS_N);
      continue;
    }


    switch(read_state) {
      case SGLF_RS_TILEID: tileid_s += (char)ch; break;
      case SGLF_RS_HASH: hash_s += (char)ch; break;
      case SGLF_RS_SEQ: seq_s += (char)ch; break;
      default: break;
    }

    line += ch;

  }

  if ((tileid_s.size()>0) || (seq_s.size()>0)) {

    if (glft_opt.use_reported_hash) {
      ts = hash_s;
    } else {
      md5str(ts, seq_s);
    }

    if (manifest.size()>0) { manifest += ' '; }
    manifest += tileid_s;
    manifest += ':';
    manifest += ts;
  }

  return 0;
}

int calc_path_lib_version(std::string &libver, FILE *ifp, glft_opt_t &glft_opt) {
  int r;
  std::string manifest;

  r = create_tile_path_manifest(manifest, ifp, glft_opt);
  if (r!=0) { return r; }


  md5str(libver, manifest);

  return 0;
}

int main(int argc, char **argv) {
  FILE *ifp=stdin, *ofp=stdout;
  int ch, opt, option_index;

  int show_help_flag=1;
  glft_opt_t glft_opt;

  int action=GLFT_NOOP;

  std::string s;

  glft_opt.use_reported_hash=0;

  while ((opt=getopt_long(argc, argv, "vVhHLPmM", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "invalid option, exiting\n");
      exit(-1);
      break;

    case 'L':
      show_help_flag=0;
      action = GLFT_LIB_VER;
      break;
    case 'P':
      show_help_flag=0;
      action = GLFT_PATH_LIB_VER;
      break;

    case 'm':
      show_help_flag=0;
      action = GLFT_LIB_MANIFEST;
      break;
    case 'M':
      show_help_flag=0;
      action = GLFT_PATH_LIB_MANIFEST;
      break;

    case 'H':
      show_help_flag=0;
      glft_opt.use_reported_hash = 1;
      break;

    case 'v':
      show_help_flag=0;
      verbose_flag=1;
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

  if ((action==GLFT_NOOP) || (show_help_flag)) {
    show_help();
    exit(0);
  }

  if (action==GLFT_LIB_VER) {
  }

  else if (action==GLFT_PATH_LIB_VER) {

    calc_path_lib_version(s, ifp, glft_opt);
    printf("%s\n", s.c_str());
  }

  else if (action == GLFT_PATH_LIB_MANIFEST) {
    create_tile_path_manifest(s, ifp, glft_opt);
    printf("%s\n", s.c_str());
  }

}
