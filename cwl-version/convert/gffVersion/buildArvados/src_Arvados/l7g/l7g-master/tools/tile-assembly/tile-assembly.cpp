#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>

#include <getopt.h>
#include <string.h>

#include <bgzf.h>

#include <vector>
#include <string>

#define TILE_ASSEMBLY_TOOL "0.1.0"

#define DEFAULT_TILE_ASSEMBLY_FILE "assembly.00.hg19.fw.gz"

static struct option long_options[] = {
  {"help",                no_argument,        NULL, 'h'},
  {"version",             no_argument,        NULL, 'v'},
  {"verbose",             no_argument,        NULL, 'V'},
  {"input",               required_argument,  NULL, 'i'},
  {0,0,0,0}
};

typedef struct ta_opt_type {
  int show_header;
  int debug;
} ta_opt_t;

void show_version() {
  printf("Version: %s\n", TILE_ASSEMBLY_TOOL);
}

void show_help() {
  printf("\n");
  printf("Tile Assembly Tool\n");
  printf("\n");
  printf("usage:\n");
  printf("\n");
  //printf("  tile-assembly <action> ...\n");
  //printf("\n");
  printf("  tile-assembly tilepath [tile-assembly-file] <hex-tilepath>\n");
  printf("\n");
  printf("    show tile assembly information for tile path\n");
  printf("\n");
  printf("\n");
  printf("  tile-assembly range [tile-assembly-file] <hex-tilepath>\n");
  printf("\n");
  printf("    show tile assembly range for tile path\n");
  printf("\n");
  printf("\n");
  return;
  printf("  [-v]        version\n");
  printf("  [-V]        verbose\n");
  printf("  [-i ifn]    input\n");
  printf("\n");
}


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


int print_bgzf(const char *afn, long byte_offset, size_t byte_len) {
  int i, r;
  BGZF *bgzfp;
  char buf[1024];
  ssize_t sz=0;
  size_t buflen=1024, cur_len=0;

  if (!(bgzfp = bgzf_open(afn, "r"))) { return -1; }
  r = bgzf_index_load(bgzfp, afn, ".gzi");
  if (r<0) { goto cleanup_print_bgzf; }

  r = bgzf_useek(bgzfp, byte_offset, SEEK_SET);
  if (r<0) { goto cleanup_print_bgzf; }

  while (cur_len < byte_len) {
    if ((cur_len + buflen) <= byte_len ) {
      sz = bgzf_read(bgzfp, buf, sizeof(char)*buflen);
    } else {
      sz = bgzf_read(bgzfp, buf, sizeof(char)*(byte_len - cur_len));
    }
    if (sz<0) { r=-1; goto cleanup_print_bgzf; }
    for (i=0; i<sz; i++) { printf("%c", buf[i]); }
    cur_len += sz;
  }

  r=0;

cleanup_print_bgzf:
  bgzf_close(bgzfp);
  return r;
}

int tileassembly_index_info(const char *afn, const char *sfx, int tilepath, std::string &name, int *byte_offset, int *byte_len) {
  FILE *idx_fp;
  std::string idx_fn, line, field;
  std::vector< std::string > tok_v;
  int ch;
  const char *chp;
  char stilepath[32];

  snprintf(stilepath, 32, "%04x", tilepath);

  idx_fn = afn;
  idx_fn += sfx;
  if ((idx_fp = fopen(idx_fn.c_str(), "r"))==NULL) { return -1; }

  while (!feof(idx_fp)) {
    ch = fgetc(idx_fp);
    if ((ch=='\n') || (ch==EOF)) {
      if (line.size()==0) { continue; }

      split_tok_ch(tok_v, line, '\t');
      line.clear();

      if (tok_v.size()!=5) { continue; }

      name = tok_v[0];

      field = ":";
      field += stilepath;

      // tile path not found
      //
      if ((chp=strstr(tok_v[0].c_str(), field.c_str()))==NULL) {
        continue;
      }

      //printf("chp: %p %i %i\n", chp, (int)(tok_v[0].c_str() + tok_v[0].size() - chp), (int)field.size());

      if ((size_t)(tok_v[0].c_str() + tok_v[0].size() - chp) != field.size()) {
        continue;
      }

      *byte_offset = atoi(tok_v[2].c_str());
      *byte_len    = atoi(tok_v[1].c_str());

      fclose(idx_fp);

      return 0;
    }

    line += (char)ch;
  }

  fclose(idx_fp);

  return -1;
}

int print_tilepath(const char *afn, int tilepath, ta_opt_t *opt) {
  int r;
  FILE *idx_fp;
  std::string name;
  int ch;
  const char *chp;

  int byte_offset, byte_len;
  int print_header=0;

  print_header = opt->show_header;

  r = tileassembly_index_info(afn, ".fwi", tilepath, name, &byte_offset, &byte_len);
  if (r<0) { return r; }

  if (print_header) { printf(">%s\n", name.c_str()); }
  print_bgzf(afn, (long)byte_offset, (size_t)byte_len);
  return 0;
}


int tileassembly_end_pos(const char *afn, int tilepath, int *tilestep, int *ref_pos) {
  int i, r;
  BGZF *bgzfp;
  char buf[1024];
  ssize_t sz=0;
  size_t buflen=1024, cur_len=0;
  int byte_offset, byte_len;
  std::string name;
  std::string line;
  int char_count=0;
  int ch;
  std::vector< std::string > tok_v;

  int cur_tilestep=-1, cur_ref_pos=-1;

  r = tileassembly_index_info(afn, ".fwi", tilepath, name, &byte_offset, &byte_len);
  if (r<0) { return -1; }

  if (!(bgzfp = bgzf_open(afn, "r"))) { return -1; }
  r = bgzf_index_load(bgzfp, afn, ".gzi");
  if (r<0) { goto cleanup_print_bgzf; }

  r = bgzf_useek(bgzfp, byte_offset, SEEK_SET);
  if (r<0) { goto cleanup_print_bgzf; }

  r=-1;
  while ((ch=bgzf_getc(bgzfp))>=0) {
    char_count++;

    if ((ch<0) || (ch=='\n') || (char_count==byte_len)) {

      split_tok_ch(tok_v, line, '\t');
      line.clear();

      if (tok_v.size()!=2) { continue; }
      cur_tilestep = strtol(tok_v[0].c_str(), NULL, 16);
      cur_ref_pos = atoi(tok_v[1].c_str());
      r=0;
      continue;
    }

    if (char_count>=byte_len) { break; }

    line += (char)ch;
  }

  *tilestep = cur_tilestep+1;
  *ref_pos = cur_ref_pos;

cleanup_print_bgzf:
  bgzf_close(bgzfp);
  return r;
}

int print_range(const char *afn, int tilepath, ta_opt_t *opt) {
  int i, j, k, r;
  int prev_byte_offset, prev_byte_len;
  int byte_offset, byte_len;

  std::string tilepath_chrom, name, ref_name;
  std::vector< std::string > tok_v;

  int start_pos = 0;
  int show_header=0;

  int end_step, end_pos;

  show_header = opt->show_header;

  tileassembly_index_info(afn, ".fwi", tilepath, name, &byte_offset, &byte_len);
  split_tok_ch(tok_v, name, ':');
  if (tok_v.size()<3) { return -1; }

  ref_name = tok_v[0];
  tilepath_chrom = tok_v[1];

  if (tilepath>0) {
    tileassembly_index_info(afn, ".fwi", tilepath-1, name, &byte_offset, &byte_len);

    split_tok_ch(tok_v, name, ':');
    if (tok_v.size()<3) { return -1; }

    if (tok_v[1] != tilepath_chrom) { start_pos = 0; }
    else {
      r = tileassembly_end_pos(afn, tilepath-1, &k, &start_pos);
    }

  }

  r = tileassembly_end_pos(afn, tilepath, &end_step, &end_pos);

  if (show_header) {
    printf("#nstep\tbeg\tend\tchrom_name\tref_name\n");
  }
  printf("%i\t%i\t%i\t%s\t%s\n",
      end_step,
      start_pos,
      end_pos,
      tilepath_chrom.c_str(),
      ref_name.c_str() );

  return 0;
}

int main(int argc, char **argv ) {
  int i, j, k;
  int ch;
  int option_index;
  int verbose_flag = 0;
  int tilepath=-1;

  std::string assembly_fn, stilepath, action;

  ta_opt_t opt;

  char *chp;
  assembly_fn = DEFAULT_TILE_ASSEMBLY_FILE;

  opt.show_header = 0;
  opt.debug=0;

  chp = getenv("TILE_ASSEMBLY");
  if (chp) { assembly_fn = chp; }

  if (argc == 3) {
    action = argv[1];
    stilepath = argv[2];
  }

  else if (argc == 4) {
    action = argv[1];
    assembly_fn = argv[2];
    stilepath = argv[3];
  }

  else {
    show_help();
    exit(0);
  }

  tilepath = strtol(stilepath.c_str(), NULL, 16);

  /*
  while ((opt = getopt_long(argc-1, argv+1, "hvVi:", long_options, &option_index))!=-1) switch (opt) {
    case 0:
      fprintf(stderr, "sanity error, invalid option to parse, exiting\n");
      exit(-1);
      break;
    case 'v': show_version(); exit(0); break;
    case 'V': verbose_flag=1; break;
    case 'h':
    default: show_help(); exit(0); break;
  }
  */

  if (action == "tilepath") {
    print_tilepath(assembly_fn.c_str(), tilepath, &opt);
  }

  else if (action == "range") {
    opt.show_header=1;
    print_range(assembly_fn.c_str(), tilepath, &opt);
  }

  else {
    show_help();
    exit(0);
  }

}
