#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include <string>
#include <map>
#include <vector>

typedef struct sglf_type {
  std::string tileid;
  std::string md5str;
  std::string seq;
} sglf_t;

int parse_tileid(int &tilepath, int &tilever, int &tilestep, int &tilevar, int &tilespan, std::string &tileid_str) {
  int i, state = 0;
  size_t pos=0, n;
  long int li;

  n = tileid_str.size();

  li = strtol(tileid_str.c_str() + pos, NULL, 16);
  if ((errno == ERANGE) || (errno == EINVAL)) { return -1; }
  tilepath = (int)li;

  for (; (pos<n) && (tileid_str[pos] != '.'); pos++);
  pos++;
  if (pos>=n) { return -2; }

  //--

  li = strtol(tileid_str.c_str() + pos, NULL, 16);
  if ((errno == ERANGE) || (errno == EINVAL)) { return -3; }
  tilever= (int)li;

  for (; (pos<n) && (tileid_str[pos] != '.'); pos++);
  pos++;
  if (pos>=n) { return -4; }

  //--

  li = strtol(tileid_str.c_str() + pos, NULL, 16);
  if ((errno == ERANGE) || (errno == EINVAL)) { return -5; }
  tilestep= (int)li;

  for (; (pos<n) && (tileid_str[pos] != '.'); pos++);
  pos++;
  if (pos>=n) { return -6; }

  //--

  li = strtol(tileid_str.c_str() + pos, NULL, 16);
  if ((errno == ERANGE) || (errno == EINVAL)) { return -7; }
  tilevar = (int)li;

  for (; (pos<n) && (tileid_str[pos] != '+'); pos++);
  pos++;
  if (pos>=n) { return -8; }

  //--

  li = strtol(tileid_str.c_str() + pos, NULL, 16);
  if ((errno == ERANGE) || (errno == EINVAL)) { return -7; }
  tilespan = (int)li;

  return 0;
}

int read_sglf_line(FILE *fp,
    int &tilepath, int &tilever, int &tilestep, int &tilevar, int &tilespan,
    std::string &m5str,
    std::string &seq) {
  int ch, state=0, char_count=0;
  std::string buf;

  // skip past whitespace
  //
  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch!=' ') && (ch!='\n') && (ch!=EOF)) {
      ungetc(ch, fp);
      break;
    }
  }

  if (ch==EOF) { return 2; }

  buf.clear();
  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch=='\n') || (ch==EOF)) { break; }

    char_count++;

    if (ch==',') {

      if      (state == 0) { parse_tileid(tilepath, tilever, tilestep, tilevar, tilespan, buf); }
      else if (state == 1) { m5str = buf; }

      buf.clear();
      state++;
      continue;
    }

    buf += (char)ch;

  }

  seq = buf;

  if (char_count>0) {
    if (state != 2) { return -1; }
    if (ch==EOF) { return 1; }
  }
  else if (char_count==0) {
    if (ch==EOF) { return 2; }
  }

  return 0;
}

typedef struct tile_type {
  int path, ver, step, var, span;
  //std::string md5str;
  //std::string seq;
} tile_t;

int sglf_merge_and_print(FILE *ofp, FILE *src_fp, FILE *add_fp) {
  int r, src_line_no=0, add_line_no=0;
  int varid=-1;

  tile_t src, src_prev, add, add_prev;

  std::string src_buf, add_buf;
  std::string src_m5, src_seq;
  std::string add_m5, add_seq;

  int src_print_prev=0, add_print_prev=0, add_beg=1;

  std::map< std::string, int > src_m5_map;
  std::map< std::string, int >::iterator srch;

  add.path=-1;
  src.path=-1;

  while ( (!feof(src_fp)) && (!feof(add_fp)) ) {

    varid = -1;
    src_m5_map.clear();

    if (src_print_prev) {
      src_m5_map[src_m5] = 1;
      printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
        src.path, src.ver, src.step, src.var, src.span,
        src_m5.c_str(),
        src_seq.c_str());
    }

    while (!feof(src_fp)) {
      r = read_sglf_line(src_fp,
                         src.path, src.ver, src.step, src.var, src.span,
                         src_m5,
                         src_seq);
      src_line_no++;
      if (r<0) { fprintf(stderr, "ERROR on source sglf line %i\n", src_line_no); continue; }
      else if (r==1) {
        // we've read a line but got EOF instead of newline at the end, let the logic outside
        // of this loop print the result.
        src_print_prev=1;
        continue;
      }
      else if (r==2) {

        // non-error eof
        //
        src_print_prev=0;
        continue;
      }

      if (src_print_prev==0) {
        src_prev.path = src.path;
        src_prev.ver  = src.ver;
        src_prev.step = src.step;
        src_prev.var  = src.var;
      }

      src_print_prev=1;

      if ((src.path != src_prev.path) ||
          (src.ver  != src_prev.ver) ||
          (src.step != src_prev.step)) {

        src_prev.path = src.path;
        src_prev.ver  = src.ver;
        src_prev.step = src.step;

        break;
      }

      src_m5_map[src_m5] = 1;
      printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
          src.path, src.ver, src.step, src.var, src.span,
          src_m5.c_str(),
          src_seq.c_str());
      src_prev.path = src.path;
      src_prev.ver  = src.ver;
      src_prev.step = src.step;

      varid = src.var;

      src_prev.path = src.path;
      src_prev.ver  = src.ver;
      src_prev.step = src.step;

    }

    if ((add.path == src.path) &&
        (add.step >= src.step)) {
      continue;
    }

    if (add_print_prev) {

      if ( (add.path < src.path) ||
           ((add.path == src.path) && (add.step < src.step)) ) {

        srch = src_m5_map.find(add_m5);
        if (srch == src_m5_map.end()) {

          varid++;

          printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
              //add.path, add.ver, add.step, add.var, add.span,
              add.path, add.ver, add.step, varid, add.span,
              add_m5.c_str(),
              add_seq.c_str());
        }
        else {
          //printf("# add skip %s\n", add_m5.c_str());
        }

      }

    }

    add_beg=1;
    while (!feof(add_fp)) {
      r = read_sglf_line(add_fp,
                         add.path, add.ver, add.step, add.var, add.span,
                         add_m5,
                         add_seq);

      add_line_no++;
      if (r<0) { fprintf(stderr, "ERROR on other sglf line %i\n", add_line_no); continue; }
      else if (r==1) {
        // we've read a line but got EOF instead of newline at the end, let the logic outside
        // of this loop print the result.
        add_print_prev=1;
        continue;
      }
      else if (r==2) {

        // non-error eof
        //
        add_print_prev=0;
        continue;
      }

      if (add_beg==1) {
        add_prev.path = add.path;
        add_prev.ver  = add.ver;
        add_prev.step = add.step;
      }
      add_beg=0;


      add_print_prev = 1;

      if ( (add.path < src.path) ||
           ((add.path == src.path) && (add.step < src.step)) ) {

        srch = src_m5_map.find(add_m5);
        if (srch == src_m5_map.end()) {

          varid++;
          printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
              add.path, add.ver, add.step, varid, add.span,
              add_m5.c_str(),
              add_seq.c_str());
        }
        else {
          //printf("# skipping %s\n", add_m5.c_str());
        }

        if (add.step != add_prev.step) { varid=-1; }

      }
      else {
        break;
      }

      add_prev.path = add.path;
      add_prev.ver  = add.ver;
      add_prev.step = add.step;

    }

  }

  // Doa final process of the src sglf
  //
  if (src_print_prev) {
    src_m5_map[src_m5] = 1;
    printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
      src.path, src.ver, src.step, src.var, src.span,
      src_m5.c_str(),
      src_seq.c_str());
  }


  //DEBUG
  //printf("...cp\n"); fflush(stdout);

  while (!feof(src_fp)) {
    r = read_sglf_line(src_fp,
                       src.path, src.ver, src.step, src.var, src.span,
                       src_m5,
                       src_seq);
    src_line_no++;
    if (r<0) { fprintf(stderr, "ERROR on source sglf line %i\n", src_line_no); continue; }
    else if (r==2) { continue; }

    /*
    if (src_print_prev==0) {
      src_prev.path = src.path;
      src_prev.ver  = src.ver;
      src_prev.step = src.step;
      src_prev.var  = src.var;
    }

    src_print_prev=1;

    if ((src.path != src_prev.path) ||
        (src.ver  != src_prev.ver) ||
        (src.step != src_prev.step)) {

      src_prev.path = src.path;
      src_prev.ver  = src.ver;
      src_prev.step = src.step;

      break;
    }
    */

  //DEBUG
  //printf("...cp1\n"); fflush(stdout);

    src_m5_map[src_m5] = 1;
    printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
        src.path, src.ver, src.step, src.var, src.span,
        src_m5.c_str(),
        src_seq.c_str());
    src_prev.path = src.path;
    src_prev.ver  = src.ver;
    src_prev.step = src.step;

    varid = src.var;

    src_prev.path = src.path;
    src_prev.ver  = src.ver;
    src_prev.step = src.step;

  }

  // Do a final process on the 'add' sglf
  //
  if (add_print_prev) {

    srch = src_m5_map.find(add_m5);
    if (srch == src_m5_map.end()) {

      varid++;

      printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
          add.path, add.ver, add.step, varid, add.span,
          add_m5.c_str(),
          add_seq.c_str());
    }
    else {
      //printf("# add skip (fin) %s\n", add_m5.c_str());
    }

  }

  while (!feof(add_fp)) {
    r = read_sglf_line(add_fp,
                       add.path, add.ver, add.step, add.var, add.span,
                       add_m5,
                       add_seq);
    add_line_no++;
    if (r<0) { fprintf(stderr, "ERROR on other sglf line %i\n", add_line_no); continue; }
    else if (r==2) { continue; }
    //else if (r==1) { continue; }

    add_print_prev = 1;

    srch = src_m5_map.find(add_m5);
    if (srch == src_m5_map.end()) {

      varid++;
      printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
          add.path, add.ver, add.step, varid, add.span,
          add_m5.c_str(),
          add_seq.c_str());
    }
    else {
      //printf("# skipping %s\n", add_m5.c_str());
    }

  }



}

void show_help(void) {
  printf("\n");
  printf("usage: merge-sglf <source-sglf> <new-sglf>\n");
  printf("\n");
}

int main(int argc, char **argv) {
  FILE *src_fp, *add_fp;

  if (argc!=3) {
    show_help();
    exit(-1);
  }

  src_fp = fopen(argv[1], "r");
  if (!src_fp) {
    perror(argv[1]);
    exit(-1);
  }

  add_fp = fopen(argv[2], "r");
  if (!add_fp) {
    perror(argv[2]);
    exit(-2);
  }

  sglf_merge_and_print(stdout, src_fp, add_fp);

  fclose(src_fp);
  fclose(add_fp);

}
