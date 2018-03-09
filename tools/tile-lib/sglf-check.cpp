#include <stdio.h>
#include <stdlib.h>

#include <string.h>
#include <errno.h>

#include <openssl/md5.h>

#include <string>
#include <map>
#include <vector>

std::string errmsg;

enum SGLF_READSTATE_ENUM {
  SGLF_RS_ERR = -1,
  SGLF_RS_OK = 0,
  SGLF_RS_OK_EOF = 1,
  SGLF_RS_EOF = 2,
};

// Helper function to create an ASCII representation
// of the MD5 digest from the sequence `seq`
//
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
  if ((errno == ERANGE) || (errno == EINVAL)) { return -8; }
  tilespan = (int)li;

  for (;
      (pos<n) &&
      ( ((tileid_str[pos] >= '0') && (tileid_str[pos] <= '9')) ||
        ((tileid_str[pos] >= 'a') && (tileid_str[pos] <= 'f')) ||
        ((tileid_str[pos] >= 'A') && (tileid_str[pos] <= 'F')) ) ;
        pos++);

  if (pos!=n) { return -9; }

  return 0;
}


int read_sglf_line(FILE *fp,
    int &tilepath, int &tilever, int &tilestep, int &tilevar, int &tilespan,
    std::string &m5str,
    std::string &seq) {
  int ret=0;
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

      if      (state == 0) {
        ret=parse_tileid(tilepath, tilever, tilestep, tilevar, tilespan, buf);
        if (ret<0) {
          errmsg = "invalid tileid: " + buf;
          return SGLF_RS_ERR;
        }
      }
      else if (state == 1) { m5str = buf; }

      buf.clear();
      state++;
      continue;
    }

    buf += (char)ch;

  }

  seq = buf;

  if (char_count>0) {
    if (state != 2) {
      errmsg = "invalid state";
      return SGLF_RS_ERR;
    }
    if (ch==EOF) { return SGLF_RS_OK_EOF; }
  }
  else if (char_count==0) {
    if (ch==EOF) { return SGLF_RS_EOF; }
  }

  return SGLF_RS_OK;
}

int main(int argc, char **argv) {
  int r;
  std::string fn;
  FILE *fp=stdin;

  int i, line_no=0;

  int tilepath, tilever, tilestep, tilevar, tilespan;
  int prv_tilepath=-1, prv_tilever=-1, prv_tilestep=-1, prv_tilevar=-1, prv_tilespan=-1;

  std::string m5str, seq, check_m5str;

  int opt_check_hash = 1;



  if (argc>1) { fn = argv[1]; }
  else {
    printf("usage:\n\n    sglf-check <sglf-stream>\n\n");
    exit(0);
  }

  if (fn != "-") {
    fp = fopen(fn.c_str(), "r");
    if (!fp) {
      perror(fn.c_str());
      exit(errno);
    }
  }

  while (!feof(fp)) {
    m5str.clear();
    seq.clear();
    r = read_sglf_line(fp, tilepath, tilever, tilestep, tilevar, tilespan, m5str, seq);

    line_no++;

    if (r==SGLF_RS_ERR) {
      printf("ERROR: got %i at line_no %i (%s)\n", r, line_no, errmsg.c_str());
      exit(-1);
    }
    else if (r==SGLF_RS_EOF) { continue; }

    if (tilespan <= 0) {
      printf("ERROR: got %i for tile span at line_no %i (%04x.%02x.%04x.%03x+%x)\n",
          tilespan,
          line_no,
          tilepath, tilever, tilestep, tilevar, tilespan);
      exit(-2);
    }

    if (opt_check_hash) {
      md5str(check_m5str, seq);
      if (check_m5str != m5str) {
        printf("ERROR: got %s (%04x.%02x.%04x.%03x+%x) for tile sequence hash but reported as %s at line_no %i\n",
            check_m5str.c_str(),
            tilepath, tilever, tilestep, tilevar, tilespan,
            m5str.c_str(),
            line_no);
        exit(-6);
      }
    }

    for (i=0; i<seq.size(); i++) {
      if ( (seq[i]!='a') &&
           (seq[i]!='c') &&
           (seq[i]!='g') &&
           (seq[i]!='t') ) {
        printf("ERROR: sequence contains not 'a', 'c', 'g', 't' characters (seq pos %i) at line_no %i\n",
            i, line_no);
        exit(-7);
      }
    }

    if (tilepath != prv_tilepath) { }
    else if (prv_tilever != tilever) { }
    else if (prv_tilestep != tilestep) {
      if ((tilestep - prv_tilestep) < 0) {
        printf("ERROR: tilestep jump non increasing (%04x.%02x.%04x.%03x+%x to %04x.%02x.%04x.%03x+%x) at line_no %i\n",
            tilepath, tilever, tilestep, tilevar, tilespan,
            prv_tilepath, prv_tilever, prv_tilestep, prv_tilevar, prv_tilespan, line_no);
        exit(-3);
      }

      if (tilevar != 0) {
        printf("ERROR: tilevar not 0 at beginning of tilestep block (%04x.%02x.%04x.%03x+%x) at line_no %i\n",
            tilepath, tilever, tilestep, tilevar, tilespan,
            line_no);
        exit(-4);
      }
    }
    else {

      if ((tilevar - prv_tilevar) != 1) {
        printf("ERROR: tilevar jump not 1 (%04x.%02x.%04x.%03x+%x to %04x.%02x.%04x.%03x+%x) at line_no %i\n",
            tilepath, tilever, tilestep, tilevar, tilespan,
            prv_tilepath, prv_tilever, prv_tilestep, prv_tilevar, prv_tilespan, line_no);
        exit(-5);
      }

    }

    prv_tilepath = tilepath;
    prv_tilever  = tilever;
    prv_tilestep = tilestep;
    prv_tilevar  = tilevar;
    prv_tilespan = tilespan;
  }


  if (fp!=stdin) { fclose(fp); }

  printf("ok\n");

  exit(0);
}
