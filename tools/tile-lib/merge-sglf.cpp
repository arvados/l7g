#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include <string>
#include <map>
#include <vector>

#define SGLF_MERGE_VERSION "0.1.4"

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

enum SGLF_READSTATE_ENUM {
  SGLF_RS_ERR = -1,
  SGLF_RS_OK = 0,
  SGLF_RS_OK_EOF = 1,
  SGLF_RS_EOF = 2,
};

// Read in an SGLF line and store the tilepath, tile library version, tilestep and tile span into the
// appropriate variables.
// Fill the `m5st` and `seq` variables with the rad in hash and sequence from the file.
//
// Return:
//  -1 - error
//   0 - success
//   1 - EOF (line read in successfully but encountered an EOF at the end)
//   2 - EOF (no data read)
//
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

  if (ch==EOF) { return SGLF_RS_EOF; }

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
    if (state != 2) { return SGLF_RS_ERR; }
    if (ch==EOF) { return SGLF_RS_OK_EOF; }
  }
  else if (char_count==0) {
    if (ch==EOF) { return SGLF_RS_EOF; }
  }

  return SGLF_RS_OK;
}

typedef struct tile_type {
  int path, ver, step, var, span;
} tile_t;

// This does a 'zipper'-like merge of both SGLF streams.
// A whole tilestep is read from the first stream and printed
// and the second stream is read for a whole tilestep, only
// printing the tiles that haven't already been printed.
//
//
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

    // reset the variable id and map of seen seqeunce hashes
    // per tilestep block
    //
    varid = -1;
    src_m5_map.clear();


    // print the already read in line of the beginning of the
    // 'tilestep block'.
    // We want to skip on the first pass so we've used a flag.
    //
    if (src_print_prev) {
      src_m5_map[src_m5] = 1;
      printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
        src.path, src.ver, src.step, src.var, src.span,
        src_m5.c_str(),
        src_seq.c_str());
      varid = src.var;
    }

    // Read a tilestep block from the 'src' SGLF stream
    //
    while (!feof(src_fp)) {
      r = read_sglf_line(src_fp,
                         src.path, src.ver, src.step, src.var, src.span,
                         src_m5,
                         src_seq);
      src_line_no++;
      if (r<0) { fprintf(stderr, "ERROR on source sglf line %i\n", src_line_no); continue; }

      // we've read a line but got EOF instead of newline at the end, let the logic outside
      // of this loop print the result.
      //
      else if (r==SGLF_RS_OK_EOF) {
        src_print_prev=1;
        continue;
      }
      else if (r==SGLF_RS_EOF) {
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

      // We've encountered the next tilestep block, break
      //
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

    // The 'src' SGLF stream tilestep block is catching up
    // to the 'add' SGLF stream, so keep on processing
    // the 'src' SGLF stream.
    //
    if ((add.path == src.path) &&
        (add.step >= src.step)) {
      continue;
    }

    // src.step holds the beginning of the 'queued' tilestep block
    // from the 'src' stream.
    // Print the beginning of the tilestep block if appropriate
    // We want to skip on the first pass so we've used a flag.
    //
    if (add_print_prev) {

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
        else { }

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

      // we've read a line but got EOF instead of newline at the end, let the logic outside
      // of this loop print the result.
      //
      else if (r==SGLF_RS_OK_EOF) {
        add_print_prev=1;
        continue;
      }
      else if (r==SGLF_RS_EOF) {
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

          // The 'src' stream skipped a tilestep but the 'add' stream
          // hasn't, so reset the varid counter
          //
          if (add.step != add_prev.step) { varid=-1; }

          varid++;
          printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
              add.path, add.ver, add.step, varid, add.span,
              add_m5.c_str(),
              add_seq.c_str());
        }
        else { }

      }
      else { break; }

      add_prev.path = add.path;
      add_prev.ver  = add.ver;
      add_prev.step = add.step;

    }

  }

  // Do a final process of the src sglf
  //
  if (src_print_prev) {
    src_m5_map[src_m5] = 1;
    printf("%04x.%02x.%04x.%03x+%x,%s,%s\n",
      src.path, src.ver, src.step, src.var, src.span,
      src_m5.c_str(),
      src_seq.c_str());
  }

  // One of the two streams is at EOF so run through and
  // process and print both with the understadning that only
  // one will actually be processed.
  //

  while (!feof(src_fp)) {
    r = read_sglf_line(src_fp,
                       src.path, src.ver, src.step, src.var, src.span,
                       src_m5,
                       src_seq);
    src_line_no++;
    if (r<0) { fprintf(stderr, "ERROR on source sglf line %i\n", src_line_no); continue; }
    else if (r==SGLF_RS_EOF) { continue; }

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
    else { }

  }

  while (!feof(add_fp)) {
    r = read_sglf_line(add_fp,
                       add.path, add.ver, add.step, add.var, add.span,
                       add_m5,
                       add_seq);
    add_line_no++;
    if (r<0) { fprintf(stderr, "ERROR on other sglf line %i\n", add_line_no); continue; }
    else if (r==SGLF_RS_EOF) { continue; }
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
    else { }

  }

  return 0;
}

void show_help(void) {
  printf("\n");
  printf("sglf-merge version: %s\n\n", SGLF_MERGE_VERSION);
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
