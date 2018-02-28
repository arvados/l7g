#include "sglf.hpp"

//#define _SGLF_LIB_DEBUG

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

      if (field_state==2) {

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

#ifdef _SGLF_LIB_DEBUG
        printf("## %x.%x.%x.%x+%x\n", tilepath, tilever, tilestep, tilevar, span);
#endif

      }
      else if (field_state==1) {
        m5str = buf;

#ifdef _SGLF_LIB_DEBUG
        printf("## m5: %s\n", m5str.c_str());
#endif


      }
      else { return -5; }

      field_state++;

#ifdef _SGLF_LIB_DEBUG
      printf("## field_state %i\n", field_state);
#endif

      buf.clear();
      continue;
    }

    buf += (char)ch;
  }

  return 0;
}

//int sglf_get(std::string &seq, int &span, int tilepath, int tilever, int tilestep, int tilevar, sglf_t &sglf) {
//  return 0;
//}

int sglf_print(sglf_t &sglf, int tilepath, int tilever) {
  MD5_CTX ctx;
  unsigned char digest[MD5_DIGEST_LENGTH];
  int i, step, tilevar;
  size_t n_tilepath, n_tilestep, n_tilevar;
  std::string seq;

  for (i=0;  i<MD5_DIGEST_LENGTH; i++) { digest[i] = 0xff; }

  n_tilepath = ((sglf.type == SGLF_TYPE_SEQ) ? sglf.seq.size() : sglf.seq2bit.size());

  //if (tilepath >= sglf.seq.size()) { return -1; }
  if (tilepath >= (int)n_tilepath) { return -1; }
  if (tilepath < 0) { return -1; }

  n_tilestep = ((sglf.type == SGLF_TYPE_SEQ) ? sglf.seq[tilepath].size() : sglf.seq2bit[tilepath].size());

  //for (step=0; step<sglf.seq[tilepath].size(); step++) {
  for (step=0; step<n_tilestep; step++) {

    n_tilevar = ((sglf.type == SGLF_TYPE_SEQ) ? sglf.seq[tilepath][step].size() : sglf.seq2bit[tilepath][step].size());

    //for (tilevar=0; tilevar<sglf.seq[tilepath][step].size(); tilevar++) {
    for (tilevar=0; tilevar<(int)n_tilevar; tilevar++) {

      if (sglf.type == SGLF_TYPE_SEQ) {
        MD5((const unsigned char *)(sglf.seq[tilepath][step][tilevar].c_str()),
            (unsigned long)sglf.seq[tilepath][step][tilevar].size(),
            digest);
      } else {
        sglf.seq2bit[tilepath][step][tilevar].twoBitToDnaSeq(seq);
        MD5((const unsigned char *)(seq.c_str()),
            (unsigned long)(seq.size()),
            digest);
      }

      printf("%04x.%02x.%04x.%03x+%x,", tilepath, tilever, step, tilevar, sglf.span[tilepath][step][tilevar]);
      for (i=0; i<MD5_DIGEST_LENGTH; i++) { printf("%02x", digest[i]); }

      if (sglf.type == SGLF_TYPE_SEQ) {
        printf(",%s\n", sglf.seq[tilepath][step][tilevar].c_str());
      } else {
        sglf.seq2bit[tilepath][step][tilevar].twoBitToDnaSeq(seq);
        printf(",%s\n", seq.c_str());
      }

    }
  }

  return 0;
}

//------

int sglf_read_2bit(FILE *ifp, sglf_t &sglf) {
  int ch;
  unsigned long long int line_no=0, char_no=0;
  int tilepath=0, tilestep=0, tilevar=0, span=-1;
  int prev_tilepath=0, prev_tilestep=0, prev_tilevar=0, prev_span=-1;
  int tilever=-1;
  int field_state=0, r;
  std::string buf, m5str;

  std::vector< std::vector< twoBit_t > > empty_vvs;
  std::vector< twoBit_t > empty_vs;
  twoBit_t empty_twobit;

  std::vector< std::vector< int > > empty_vvi;
  std::vector< int > empty_vi;

  twoBit_t twobit_seq;

  while (!feof(ifp)) {
    ch=fgetc(ifp);
    if ((ch=='\n') || (ch==EOF)) {
      if (ch=='\n') { char_no++; line_no++; }

      if (field_state==2) {

        while (tilepath >= sglf.seq2bit.size()) {
          sglf.seq2bit.push_back(empty_vvs);
          sglf.span.push_back(empty_vvi);
        }

        while (tilestep >= sglf.seq2bit[tilepath].size()) {
          sglf.seq2bit[tilepath].push_back(empty_vs);
          sglf.span[tilepath].push_back(empty_vi);
        }

        while (tilevar >= sglf.seq2bit[tilepath][tilestep].size()) {
          sglf.seq2bit[tilepath][tilestep].push_back(empty_twobit);
          sglf.span[tilepath][tilestep].push_back(0);
        }

        sglf.seq2bit[tilepath][tilestep][tilevar].twoBitFromDnaSeq(buf.c_str());
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

#ifdef _SGLF_LIB_DEBUG
        printf("## %x.%x.%x.%x+%x\n", tilepath, tilever, tilestep, tilevar, span);
#endif

      }
      else if (field_state==1) {
        m5str = buf;

#ifdef _SGLF_LIB_DEBUG
        printf("## m5: %s\n", m5str.c_str());
#endif


      }
      else { return -5; }

      field_state++;

#ifdef _SGLF_LIB_DEBUG
      printf("## field_state %i\n", field_state);
#endif

      buf.clear();
      continue;
    }

    buf += (char)ch;
  }

  return 0;
}


#ifdef _SGLF_LIB_DEBUG

int main(int argc, char **argv) {
  int i, j, k;
  sglf_t sglf;
  int r;

  int n_tilepath=863, tilepath;
  tilepath=0;

  FILE *sglf_fp=stdin;

  r = sglf_read(sglf_fp, sglf);
  if (r<0) { printf("ERROR: sglf_read: %i\n", r); exit(-1); }

  printf(".... %i\n", (int)sglf.seq.size()); fflush(stdout);

  for (tilepath=0; tilepath<n_tilepath; tilepath++) {
    if (sglf.seq[tilepath].size() > 0) {

      printf("## printing tilepath %i (%x)\n", tilepath, tilepath);

      r = sglf_print(sglf, tilepath, 0);
      if (r<0) { printf("ERROR: sglf_print: %i\n", r); exit(-2); }
    }
  }


}

#endif
