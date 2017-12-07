#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>

#include <string>
#include <vector>

int vcf_bounds(std::string &chrom, int &start1, int &end1_inc, std::vector< std::string > &field_v) {
  long int v;
  int ref_len=0;
  int end_a_inc, end_b_inc;
  const char *chp=NULL;

  const char *x, *y;

  chrom = field_v[0];
  v=strtol(field_v[1].c_str(), NULL, 10);
  if ((errno==EINVAL) || (errno==ERANGE)) { return -1; }

  start1 = v;
  end_a_inc = start1;
  end_b_inc = start1;

  ref_len = (int)field_v[3].size();
  end_a_inc += ref_len-1;

  if (field_v[8].size() > 4) {
    if ((field_v[7][0] == 'E') &&
        (field_v[7][1] == 'N') &&
        (field_v[7][2] == 'D') &&
        (field_v[7][3] == '=') ) {
      v = strtol(field_v[8].c_str()+4, NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -1; }
      end_b_inc=v;
    }
    else if (chp = strstr(field_v[7].c_str(), ";END=")) {

  printf("v %i\n", (int)v);

      v = strtol((const char *)(chp+5), NULL, 10);
      if ((errno==EINVAL) || (errno==ERANGE)) { return -1; }
      end_b_inc=v;
    }
  }

  end1_inc = ( (end_a_inc > end_b_inc) ? end_a_inc : end_b_inc );

  return 0;
}

int main(int argc, char **argv) {
  int i, j, k, r, v;
  int ch;

  int line_no=0, line_char_count=0;
  int header=1;
  int prev_s_pos=0, field_pos=0, field_count=10;
  int prev_start1=0, prev_len=0;
  int start1, end1_inc, pos_len;
  std::string chrom, buf, field;
  FILE *ifp=stdin;

  int verbose=0;

  std::vector< std::string > field_v;
  std::string empty_s;

  for (i=0; i<field_count; i++) { field_v.push_back(empty_s); }

  while (!feof(ifp)) {
    ch=fgetc(ifp);
    if ((ch=='\n') || (ch==EOF)) {
      line_no++;

      if (line_char_count==0) {
        if (ch=='\n') { printf("\n"); }
        continue;
      }
      line_char_count=0;

      if (header) {
        if (buf[0]=='#') {
          printf("%s\n", buf.c_str());
          buf.clear();
          continue;
        }
        else {
          printf("ERROR: line %i: expoected header but no beginning '#' found\n", line_no);
          exit(-1);
        }
      }

      if (field_pos < field_v.size()) {
        field_v[field_pos] = field;
      }
      field_pos++;
      field.clear();

      if (field_pos != field_count) {
        printf("ERROR: line %i: field count mismatch, got %i, expected %i\n",
            line_no, field_pos+1, field_count);
        exit(-1);
      }

      r = vcf_bounds(chrom, start1, end1_inc, field_v);
      if (r<0) {
        printf("ERROR: line %i: could not parse line\n", line_no);
        exit(-2);
      }

      if (start1 >= (prev_start1 + prev_len)) {
        for (i=0; i<field_count; i++) {
          if (i>0) { printf("\t"); }
          printf("%s", field_v[i].c_str());
        }
        printf("\n");

        prev_start1 = start1;
        prev_len = end1_inc - start1 + 1;
      }
      else if (verbose) {
        printf("## SKIPPING %s:%i-%i\n", chrom.c_str(), start1, end1_inc);
      }

      field_pos=0;
      field.clear();
      for (i=0; i<10; i++) { field_v[i].clear(); }
      continue;
    }

    line_char_count++;

    if (header) {

      if ((line_char_count==1) && (ch!='#')) {
        header=0;
      }
      else {
        buf += (char)ch;
        continue;
      }
    }

    if (ch=='\t') {
      if (field_pos < field_v.size()) {
        field_v[field_pos] = field;
      }
      field_pos++;
      field.clear();
      continue;
    }

    field += (char)ch;

  }
}
