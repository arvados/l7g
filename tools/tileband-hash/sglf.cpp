#include "sglf.hpp"

int parse_etileid(int *tilepath, int *tilever, int *tilestep, int *tilevar, int *span, std::string &field) {
  int i, j, k, n;
  std::vector<int> comma_pos;
  int plus_pos=-1;
  long int v=-1;

  for (i=0; i<field.size(); i++) {
    if (field[i]=='.') {
      comma_pos.push_back(i);
    }
    if (field[i]=='+') {
      if (plus_pos>=0) { return -1; }
      plus_pos=i;
    }
  }

  //DEBUG
  //printf("comma_pos.size() %i, plus_pos %i\n",
  //    (int)comma_pos.size(), plus_pos);

  if (comma_pos.size()!=3) { return -2; }

  v = strtol(field.c_str(), NULL, 16);
  if ((v==0) && (errno==EINVAL)) { return -3; }
  *tilepath = (int)v;

  v = strtol(field.c_str() + (size_t)comma_pos[0]+1, NULL, 16);
  if ((v==0) && (errno==EINVAL)) { return -4; }
  *tilever = (int)v;

  v = strtol(field.c_str() + (size_t)comma_pos[1]+1, NULL, 16);
  if ((v==0) && (errno==EINVAL)) { return -5; }
  *tilestep = (int)v;

  v = strtol(field.c_str() + (size_t)comma_pos[2]+1, NULL, 16);
  if ((v==0) && (errno==EINVAL)) { return -6; }
  *tilevar = (int)v;

  v = strtol(field.c_str() + (size_t)plus_pos+1, NULL, 16);
  if ((v==0) && (errno==EINVAL)) { return -7; }
  *span = (int)v;

  return 0;
}

int sglf_read(FILE *ifp, sglf_t &sglf) {
  int ch;
  unsigned long long int line_no=0, char_no=0;
  int tilepath=0, tilestep=0, tilevar=0, span=-1;
  int prev_tilepath=0, prev_tilestep=0, prev_tilevar=0, prev_span=-1;
  int tilever=-1;
  int field_state=0, r;
  std::string buf, m5str;

  std::vector< std::vector< std::string > > empty_vvs;
  std::vector< std::string > empty_vs;
  std::string empty_str;

  std::vector< std::vector< int > > empty_vvi;
  std::vector< int > empty_vi;

  while (!feof(ifp)) {
    ch=fgetc(ifp);
    if ((ch=='\n') || (ch==EOF)) {
      if (ch=='\n') { char_no++; line_no++; }

      if (buf.size()>0) {

        while (tilepath >= sglf.seq.size()) {
          sglf.seq.push_back(empty_vvs);
          sglf.span.push_back(empty_vvi);
        }

        while (tilestep >= sglf.seq[tilepath].size()) {
          sglf.seq[tilepath].push_back(empty_vs);
          sglf.span[tilepath].push_back(empty_vi);
        }

        while (tilevar >= sglf.seq[tilepath][tilestep].size()) {
          sglf.seq[tilepath][tilestep].push_back(empty_str);
          sglf.span[tilepath][tilestep].push_back(0);
        }

        sglf.seq[tilepath][tilestep][tilevar] = buf;
        sglf.span[tilepath][tilestep][tilevar] = span;

        prev_tilepath = tilepath;
        prev_tilestep = tilestep;
        prev_tilevar = tilevar;
        prev_span = span;
      }

      buf.clear();
      field_state=0;
      continue;
    }

    char_no++;

    if (ch==',') {

      if (field_state==0) {
        r = parse_etileid(&tilepath, &tilever, &tilestep, &tilevar, &span, buf);
        if (span<1) { return -10; }
        if (r<0) { return r; }
      }
      else if (field_state==1) {
        m5str = buf;
      }
      else { return -5; }

      field_state++;
      buf.clear();
      continue;
    }

    buf += (char)ch;
  }

  return 0;
}

int sglf_get(std::string &seq, int &span, int tilepath, int tilever, int tilestep, int tilevar, sglf_t &sglf) {
  return 0;
}

int sglf_print(sglf_t &sglf, int tilepath, int tilever) {
  MD5_CTX ctx;
  unsigned char digest[MD5_DIGEST_LENGTH];
  int i;
  int step, tilevar;

  for (i=0;  i<MD5_DIGEST_LENGTH; i++) { digest[i] = 0xff; }

  if (tilepath >= sglf.seq.size()) { return -1; }
  if (tilepath < 0) { return -1; }

  for (step=0; step<sglf.seq[tilepath].size(); step++) {
    for (tilevar=0; tilevar<sglf.seq[tilepath][step].size(); tilevar++) {

      MD5((const unsigned char *)(sglf.seq[tilepath][step][tilevar].c_str()),
          (unsigned long)sglf.seq[tilepath][step][tilevar].size(),
          digest);

      printf("%04x.%02x.%04x.%03x+%x,", tilepath, tilever, step, tilevar, sglf.span[tilepath][step][tilevar]);
      for (i=0; i<MD5_DIGEST_LENGTH; i++) { printf("%02x", digest[i]); }
      printf(",%s\n", sglf.seq[tilepath][step][tilevar].c_str());

    }
  }

  return 0;
}

/*
int main(int argc, char **argv) {
  int i, j, k;
  sglf_t sglf;
  int r;

  int tilepath=0x35e;
  tilepath=0;

  FILE *sglf_fp=stdin;

  r = sglf_read(sglf_fp, sglf);
  if (r<0) { printf("ERROR: sglf_read: %i\n", r); exit(-1); }

  for (tilepath=0; tilepath<sglf.seq.size(); tilepath++) {
    if (sglf.seq[tilepath].size() > 0) {
      r = sglf_print(sglf, tilepath, 0);
      if (r<0) { printf("ERROR: sglf_print: %i\n", r); exit(-2); }
    }
  }

  //printf("... %i\n", r);
}

*/
