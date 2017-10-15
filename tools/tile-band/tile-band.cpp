#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <getopt.h>

#include <vector>
#include <string>
#include <map>

#define VERSION_STR "0.1.0"

int VERBOSE_MATCH;

typedef struct tileband_type {
  std::vector< int > v[2];
  std::vector< std::vector< int > > noc_v[2];
} tileband_t;

void tileband_knot(tileband_t &t, int knot_start, int *knot_len, int *loc_knot);

int tileband_concordance(tileband_t &a, tileband_t &b, int s, int n, int hiq_flag, int *r_match, int *r_tot) {
  int i, j, k;
  int end_noninc=0;
  int match=0, tot=0;

  int idx_a, idx_b;
  int knot_a_len, knot_b_len;
  int knot_a_loq, knot_b_loq;
  int is_match=0;

  if (s<=0) { s=0; }
  if (n<=0) { end_noninc = a.v[0].size(); }
  else { end_noninc = s+n; }
  if (end_noninc >= a.v[0].size()) { end_noninc = a.v[0].size(); }
  n = end_noninc - s;

  idx_a = s;
  idx_b = s;


  tileband_knot(a, idx_a, &knot_a_len, &knot_a_loq);
  tileband_knot(b, idx_b, &knot_b_len, &knot_b_loq);

  while ( (idx_a < (s+n)) &&
          (idx_b < (s+n)) ) {

    if (VERBOSE_MATCH) {
    printf("idx_a %i+%i (%i), idx_b %i+%i (%i) (%i / %i)\n",
        idx_a, knot_a_len, knot_a_loq,
        idx_b, knot_b_len, knot_b_loq,
       match, tot );
    }

    if (idx_a == idx_b) {

      if (hiq_flag) {

        if ((knot_a_loq==0) && (knot_b_loq==0)) {
          tot++;

          if ((knot_a_len == knot_b_len) &&
              (knot_a_loq == 0) &&
              (knot_b_loq == 0)) {
            is_match=1;
            for (i=idx_a; i<(idx_a+knot_a_len); i++) {
              if ((a.v[0][i] != b.v[0][i]) ||
                  (a.v[1][i] != b.v[1][i])) {
                is_match=0;
                break;
              }
            }
            if (is_match) {

              if (VERBOSE_MATCH) { printf("MATCH %i+%i\n", idx_a, knot_a_len); }

              match++;
            }
          }
        }

      } else {
        tot++;

        if (knot_a_len == knot_b_len) {
          is_match=1;
          for (i=idx_a; i<(idx_a+knot_a_len); i++) {
            if ((a.v[0][i] != b.v[0][i]) ||
                (a.v[1][i] != b.v[1][i])) {
              is_match=0;
              break;
            }
          }
          if (is_match) {

            if (VERBOSE_MATCH) { printf("MATCH %i+%i\n", idx_a, knot_a_len); }

            match++;
          }
        }

      }

      idx_a += knot_a_len;
      idx_b += knot_b_len;

      if ((idx_a < (s+n)) && (idx_b < (s+n))) {
        tileband_knot(a, idx_a, &knot_a_len, &knot_a_loq);
        tileband_knot(b, idx_b, &knot_b_len, &knot_b_loq);
      }
      continue;
    }

    if (idx_a < idx_b) {
      idx_a += knot_a_len;
      if (idx_a < (s+n)) {
        tileband_knot(a, idx_a, &knot_a_len, &knot_a_loq);
      }
      continue;
    }

    if (idx_a > idx_b) {
      idx_b += knot_b_len;
      if (idx_b < (s+n)) {
        tileband_knot(b, idx_b, &knot_b_len, &knot_b_loq);
      }
      continue;
    }

  }

  *r_match = match;
  *r_tot = tot;

}

void tileband_knot(tileband_t &t, int knot_start, int *knot_len, int *loc_knot) {
  int i, n;
  int kl=1;

  *loc_knot=0;

  if ((t.noc_v[0][knot_start].size()>0) || 
      (t.noc_v[1][knot_start].size()>0)) {
    *loc_knot = 1;
  }

  n = (int)t.v[0].size();
  for (i=knot_start+1; i<n; i++) {

    if ((t.v[0][i] < 0) || (t.v[1][i] < 0)) {
      kl++;

      if ((t.noc_v[0][i].size()>0) || 
          (t.noc_v[1][i].size()>0)) {
        *loc_knot = 1;
      }

      continue;
    }

    break;
  }

  *knot_len = kl;
}

void tileband_loq_print(tileband_t &t, int start, int n_orig, int fill_loq_spanning) {
  int i, j, k, val, n;
  std::vector< int > loq_band[2];

  n = (int)t.v[0].size();
  for (i=0; i<n; i++) {

    val = ( (t.noc_v[0][i].size()>0) ? -2 : t.v[0][i] );
    loq_band[0].push_back(val);

    val = ( (t.noc_v[1][i].size()>0) ? -2 : t.v[1][i] );
    loq_band[1].push_back(val);

  }


  if (fill_loq_spanning) {
    for (i=0; i<n; ) {

      if ((loq_band[0][i] > -2) && (loq_band[1][i] > -2)) {
        i++;
        continue;
      }

      k=i+1;
      while ((k<n) &&
             ((loq_band[0][k] == -1) || (loq_band[1][k] == -1))) {
        loq_band[0][k]=-2;
        loq_band[1][k]=-2;
        k++;
      }
      i++;

    }
  }



  if (start<0) { start = 0; }

  n = n_orig;
  if (n<0) {
    n=(int)loq_band[0].size();
    n-=start;
  }

  if ((start+n) > (int)(loq_band[0].size())) {
    n = (int)(loq_band[0].size()) - start;
  }

  printf("[");
  for (i=start; i<(start+n); i++) { printf(" %i", loq_band[0][i]); }
  printf("]\n");

  printf("[");
  for (i=start; i<(start+n); i++) { printf(" %i", loq_band[1][i]); }
  printf("]\n");

  return;

  /*

  if ((n<0) || (start>((int)t.v[0].size()))) { return; }

  printf("[");
  for (i=0; i<t.v[0].size(); i++) { printf(" %i", t.v[0][i]); }
  printf("]\n");

  printf("[");
  for (i=0; i<t.v[1].size(); i++) { printf(" %i", t.v[1][i]); }
  printf("]\n");

  printf("[");
  for (i=0; i<t.noc_v[0].size(); i++) {

    printf("[");
    for (j=0; j<t.noc_v[0][i].size(); j++) {
      printf(" %i", t.noc_v[0][i][j]);
    }
    printf(" ]");

  }
  printf("]\n");

  printf("[");
  for (i=0; i<t.noc_v[1].size(); i++) {

    printf("[");
    for (j=0; j<t.noc_v[1][i].size(); j++) {
      printf(" %i", t.noc_v[1][i][j]);
    }
    printf(" ]");

  }
  printf("]\n");
  */


}

void tileband_print(tileband_t &t, int start, int n) {
  int i, j, k;

  if (start<0) { start = 0; }
  if (n<0) {
    n=(int)t.v[0].size();
    n-=start;
  }

  if ((n<0) || (start>((int)t.v[0].size()))) { return; }

  printf("[");
  for (i=start; i<(start+n); i++) { printf(" %i", t.v[0][i]); }
  printf("]\n");

  printf("[");
  for (i=start; i<(start+n); i++) { printf(" %i", t.v[1][i]); }
  printf("]\n");

  printf("[");
  for (i=start; i<(start+n); i++) {

    printf("[");
    for (j=0; j<t.noc_v[0][i].size(); j++) {
      printf(" %i", t.noc_v[0][i][j]);
    }
    printf(" ]");

  }
  printf("]\n");

  printf("[");
  for (i=start; i<(start+n); i++) {

    printf("[");
    for (j=0; j<t.noc_v[1][i].size(); j++) {
      printf(" %i", t.noc_v[1][i][j]);
    }
    printf(" ]");

  }
  printf("]\n");

  return;

  /*
  printf("[");
  for (i=0; i<t.v[0].size(); i++) { printf(" %i", t.v[0][i]); }
  printf("]\n");

  printf("[");
  for (i=0; i<t.v[1].size(); i++) { printf(" %i", t.v[1][i]); }
  printf("]\n");

  printf("[");
  for (i=0; i<t.noc_v[0].size(); i++) {

    printf("[");
    for (j=0; j<t.noc_v[0][i].size(); j++) {
      printf(" %i", t.noc_v[0][i][j]);
    }
    printf(" ]");

  }
  printf("]\n");

  printf("[");
  for (i=0; i<t.noc_v[1].size(); i++) {

    printf("[");
    for (j=0; j<t.noc_v[1][i].size(); j++) {
      printf(" %i", t.noc_v[1][i][j]);
    }
    printf(" ]");

  }
  printf("]\n");
  */

}

int tileband_read(tileband_t &t, FILE *fp) {
  int ch, band_state=0, paren=0;
  std::string buf;
  std::vector< int > loq_v;

  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch==EOF) || (ch=='\n')) {
      band_state++;
      continue;
    }

    if (ch=='[') { paren++; continue; }
    if (ch==']') {
      paren--;

      if ((band_state==0) || (band_state==1)) {
        if (buf.size()>0) {
          t.v[band_state].push_back(atoi(buf.c_str()));
        }
      }

      else if ((band_state==2) || (band_state==3)) {
        if (paren==1) {
          t.noc_v[band_state-2].push_back(loq_v);
          loq_v.clear();
        }
      }

      buf.clear();
      continue;
    }

    if (ch==' ') {

      if ((band_state==0) || (band_state==1)) {
        if (buf.size()>0) {
          t.v[band_state].push_back(atoi(buf.c_str()));
        }
      }
      else if ((band_state==2) || (band_state==3)) {
        if (buf.size()>0) {
          loq_v.push_back(atoi(buf.c_str()));
        }
      }
      buf.clear();
      continue;
    }

    buf.push_back((char)ch);
  }

  return 0;
}

static struct option long_options[] = {
  {"step",          required_argument,        NULL, 's'},
  {"endstep",       required_argument,        NULL, 'S'},
  {"help",                no_argument,        NULL, 'h'},
  {"version",             no_argument,        NULL, 'v'},
  {"verbose",             no_argument,        NULL, 'V'},
  {0,0,0,0}
};

void show_version() {
  printf("%s", VERSION_STR);
}

void show_help() {
  show_version();
  printf("usage:\n");
  printf("    tile-band [-h] [-v] [-V] [-s s] [-S S] band_file\n");
  printf("  [-s s]  tile step start\n");
  printf("  [-S S]  tile step end (inclusive)\n");
  printf("  [-n n]  tile step count\n");
  printf("  [-q]    print low quality bands (-2 for low quality entry)\n");
  printf("  [-Q]    fill low quality spanning tiles\n");
  printf("  [-h]    show help (this screen)\n");
  printf("  [-v]    verbose\n");
  printf("  [-V]    version\n");
}

int main(int argc, char **argv) {
  int i, j, k;
  FILE *fp;
  tileband_t tileband_a, tileband_b;
  std::string ifn_a, ifn_b;
  int match=0, tot=0;

  int opt, option_index;
  int verbose_flag = 0;

  int start_tilestep=0, end_tilestep_inc = -1;
  int n_tilestep = -1;
  int hiq_flag = 1, fill_loq_spanning_flag=0;

  VERBOSE_MATCH=0;

  while ((opt = getopt_long(argc, argv, "hvVs:S:n:qQ", long_options, &option_index))!=-1) switch (opt) {
    case 0:
      fprintf(stderr, "invalid argument");
      exit(-1);
      break;
    case 's': start_tilestep = atoi(optarg); break;
    case 'S': end_tilestep_inc = atoi(optarg); break;
    case 'n': n_tilestep = atoi(optarg); break;
    case 'v': verbose_flag = 1; break;
    case 'V': show_version(); exit(0); break;
    case 'q': hiq_flag=0; break;
    case 'Q': fill_loq_spanning_flag=1; break;
    default:
    case 'h': show_help(); exit(0); break;
  }

  if ((argc-optind)>=1) {
    ifn_a = argv[optind];
  } else {
    show_help();
    exit(0);
  }

  if (end_tilestep_inc >= 0) {
    n_tilestep = end_tilestep_inc - start_tilestep + 1;
  }

  if (ifn_a == "-") {
    fp = stdin;
  } else if (!(fp = fopen(ifn_a.c_str(), "r"))) {
    perror(ifn_a.c_str());
    exit(-1);
  }
  tileband_read(tileband_a, fp);
  if (fp!=stdin) { fclose(fp); }

  if (!hiq_flag) {
    tileband_loq_print(tileband_a, start_tilestep, n_tilestep, fill_loq_spanning_flag);
  }
  else {
    tileband_print(tileband_a, start_tilestep, n_tilestep);
  }

}
