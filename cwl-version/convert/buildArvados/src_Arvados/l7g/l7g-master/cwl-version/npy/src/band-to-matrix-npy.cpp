#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <string.h>

#include <cstdlib>
#include <iostream>
#include <map>
#include <string>

#include <vector>
#include <string>
#include <iterator>
#include <complex>

#include "cnpy.h"

#define N_ALLELE 1

int TILEPATH;

void print_onehot(int x, int len, char const *s) {
  int i;

  for (i=0; i<len; i++) {
    if ((i>0) && s) { printf("%s", s); }
    if (i==x) { printf("1"); }
    else { printf("0"); }
  }
}

typedef struct tilepath_vec_type {
  std::string name;
  std::vector<int> allele[2];
  std::vector<int> loq_flag[2];
} tilepath_vec_t;

void save_npy_int(const char *fn, std::vector<tilepath_vec_t> &tv, char opt) {
  int i, j, n, cur;
  unsigned int shape[] = {0,0};

  int *biga;
  int val;

  n = tv[0].allele[0].size();
  biga = new int[tv.size()*2*n];


  cur = 0;
  for (i=0; i<tv.size(); i++) {
    for (j=0; j<n; j++) {

      val = tv[i].allele[0][j];

      if ((opt == 'I') || (opt == 'C')) {
        if (tv[i].loq_flag[0][j]) { val = -2; }

      }

      // interleave
      if ((opt == 'i') || (opt=='I')) {
        biga[cur++] = val;

        val = tv[i].allele[1][j];
        if (opt == 'I') { if (tv[i].loq_flag[1][j]) { val = -2; } }

        biga[cur++] = val;
      }

      //concat
      else if ((opt == 'c') || (opt == 'C')) {
        biga[cur++] = val;
      }

    }

    // second part of concat
    if ((opt=='c') || (opt=='C')) {
      for (j=0; j<n; j++) {
        val = tv[i].allele[1][j];
        if (tv[i].loq_flag[1][j]) { val = -2; }

        biga[cur++] = val;
      }
    }
  }

  shape[0] = (int)(tv.size());
  shape[1] = (int)(2*n);

  cnpy::npy_save(fn, biga, (const unsigned int *)shape, 2, "w");

  delete biga;
}

void save_npy_double(const char *fn, std::vector<tilepath_vec_t> &tv, char opt) {
  int i, j, n, cur;
  unsigned int shape[] = {0,0};

  double *biga, val;

  n = tv[0].allele[0].size();
  biga = new double[2*tv.size()*n];


  cur = 0;
  for (i=0; i<tv.size(); i++) {
    for (j=0; j<n; j++) {

      val = ( (tv[i].loq_flag[0][j]) ? NAN : ((double)tv[i].allele[0][j]) );

      // interleave
      if (opt == 'i') {
        biga[cur++] = val;
        val = ( (tv[i].loq_flag[1][j]) ? NAN : ((double)tv[i].allele[1][j]) );
        biga[cur++] = val;
      }

      //concat
      else if (opt == 'c') {
        biga[cur++] = val;
      }

    }

    // second part of concat
    if (opt=='c'){
      for (j=0; j<n; j++) {
        val = ( (tv[i].loq_flag[1][j]) ? NAN : ((double)tv[i].allele[1][j]) );
        biga[cur++] = tv[i].allele[1][j];
      }
    }
  }

  shape[0] = (int)(2*tv.size());
  shape[1] = (int)(n);

  cnpy::npy_save(fn, biga, (const unsigned int *)shape, 2, "w");

  delete biga;
}

void print_tilepath_vec(tilepath_vec_t &tv) {
  int i, j;
  printf("%s\n", tv.name.c_str());
  for (i=0; i<2; i++) {
    for (j=0; j<tv.allele[i].size(); j++) {
      printf(" %i [%i]", tv.allele[i][j], tv.loq_flag[i][j]);
    }
    printf("\n");
  }
}

void print_tilepath_vecs(std::vector<tilepath_vec_t> &tv) {
  size_t n, tilepath_n;
  int i, j, pos;
  int max_val;
  int len;

  int print_header = 1;
  int hotpos;

  n = tv.size();
  tilepath_n = tv[0].allele[0].size();

  /*
  if (print_header) {
    printf(" pos \\ name  |");
    for (i=0; i<n; i++) {
      printf(" %s", tv[i].name.c_str());
    }
    printf("\n");
    printf("-------------|--\n");
  }
  */

  for (pos=0; pos<tilepath_n; pos++) {

    max_val = 0;
    for (i=0; i<n; i++) {

      if ((tv[i].loq_flag[0][pos]==0) &&
          (max_val < tv[i].allele[0][pos])) {
        max_val = tv[i].allele[0][pos];
      }

      if ((tv[i].loq_flag[1][pos]==0) &&
          (max_val < tv[i].allele[1][pos])) {
        max_val = tv[i].allele[1][pos];
      }

    }

    len = (int)(max_val)+1;

    for (hotpos=0; hotpos<len; hotpos++) {
      for (i=0; i<n; i++) {

        if (i>0) { printf(" "); }
        else if (print_header) {
          //printf("pos%03x.%03x.u | ", pos, hotpos);
          if (TILEPATH>=0) {
            printf("%04x.%03x(%03x)u ", TILEPATH, pos, hotpos);
          } else {
            printf("pos%03x(%03x)u ", pos, hotpos);
          }
        }

        if (tv[i].loq_flag[0][pos]==0) {
          printf("%i", (tv[i].allele[0][pos] == hotpos) ? 1 : 0);
        } else {
          printf("NaN");
        }

      }

      printf("\n");

    }

    for (hotpos=0; hotpos<len; hotpos++) {

      for (i=0; i<n; i++) {

        if (i>0) { printf(" "); }
        else if (print_header) {
          //printf("pos%03x.%03x.v | ", pos, hotpos);
          if (TILEPATH>=0) {
            printf("%04x.%03x(%03x)v ", TILEPATH, pos, hotpos);
          } else {
            printf("pos%03x(%03x)v ", pos, hotpos);
          }
        }

        if (tv[i].loq_flag[1][pos]==0) {
          printf("%i", (tv[i].allele[1][pos] == hotpos) ? 1 : 0);
        } else {
          printf("NaN");
        }

      }

      printf("\n");
    }

    //printf("\n");
  }

}

void spot_test() {
  int i, j, ch;
  int read_line = 0;
  int step=0;

  std::string s;

  std::vector<tilepath_vec_t> ds;
  tilepath_vec_t cur_ds;

  int pcount=0;

  cur_ds.allele[0].push_back(0);
  cur_ds.allele[1].push_back(2);

  cur_ds.loq_flag[0].push_back(1);
  cur_ds.loq_flag[1].push_back(0);

  cur_ds.allele[0].push_back(0);
  cur_ds.allele[1].push_back(0);

  cur_ds.loq_flag[0].push_back(0);
  cur_ds.loq_flag[1].push_back(0);

  cur_ds.name = "ds0";

  ds.push_back(cur_ds);

  cur_ds.allele[0].clear();
  cur_ds.allele[1].clear();
  cur_ds.loq_flag[0].clear();
  cur_ds.loq_flag[1].clear();
  cur_ds.name.clear();

  cur_ds.allele[0].push_back(0);
  cur_ds.allele[1].push_back(0);

  cur_ds.allele[0].push_back(3);
  cur_ds.allele[1].push_back(1);

  cur_ds.loq_flag[0].push_back(0);
  cur_ds.loq_flag[1].push_back(0);

  cur_ds.loq_flag[0].push_back(0);
  cur_ds.loq_flag[1].push_back(0);

  cur_ds.name = "ds1";

  ds.push_back(cur_ds);

  //print_tilepath_vec(cur_ds);

  print_tilepath_vecs(ds);
  exit(0);
}

int main(int argc, char **argv) {
  int i, j, ch;
  int read_line = 0;
  int step=0;

  std::vector<std::string> names;
  std::string s;

  std::vector<tilepath_vec_t> ds;
  tilepath_vec_t cur_ds;

  int pcount=0;
  int state_mod = 0;
  int cur_allele = 0;
  int loq_flag = 0;

  char *ofn;

  ofn = strdup("oot.npy");

  TILEPATH=-1;

  if (argc>1) {
    TILEPATH = atoi(argv[1]);
    if (argc>2) {
      free(ofn);
      ofn = strdup(argv[2]);
    }
  }



  s.clear();
  while (ch!=EOF) {
    ch = fgetc(stdin);

    if (ch==EOF) { break; }
    if (ch=='\n') {
      state_mod = (state_mod+1)%4;
      pcount=0;

      if (state_mod==0) {
        ds.push_back(cur_ds);
        cur_ds.allele[0].clear();
        cur_ds.allele[1].clear();
        cur_ds.loq_flag[0].clear();
        cur_ds.loq_flag[1].clear();
        cur_ds.name.clear();
      }

      continue;
    }
    if (ch=='[') {

      loq_flag=0;
      if (state_mod>=2) {

        pcount++;
        while (pcount>1) {
          ch = fgetc(stdin);

          if (ch==EOF) {
            printf("ERROR: premature eof\n");
            exit(1);
          }


          if (ch==']') {
            cur_allele = state_mod%2;
            cur_ds.loq_flag[cur_allele].push_back(loq_flag);
            pcount--;
            continue;
          }

          if (ch=='[') { pcount++; continue; }
          if (ch==' ') { continue; }

          loq_flag=1;
        }
      }

      continue;
    }
    if ((ch==' ') || (ch==']')) {

      if (s.size() == 0) { continue; }

      if (state_mod==0) {
        cur_ds.allele[0].push_back(atoi(s.c_str()));
      } else if (state_mod==1) {
        cur_ds.allele[1].push_back(atoi(s.c_str()));
      } else if (state_mod==2) {
        if (ch==']') { pcount--; }
      }
      s.clear();
      continue;
    }
    s += ch;

  }

  //for (i=0; i<ds.size(); i++) { print_tilepath_vec(ds[i]); }

  save_npy_int(ofn, ds, 'I');
  free(ofn);
  //save_npy_double("test0.npy", ds, 'i');

  //print_tilepath_vecs(ds);


  /*
  for (i=0; i<N_ALLELE; i++) {

    for (ch=fgetc(stdin); (ch==' ') || (ch=='\n'); ch=fgetc(stdin)) ;

    if (ch!='[') { fprintf(stderr, "bad input"); fflush(stderr); exit(1); }
    pcount = 1;

    step=0;

    while (pcount>0) {
      ch = fgetc(stdin);

      if (ch==']') {
        pcount--;
        step++;
        continue;
      }

      if (ch=='[') { pcount++; continue; }
      if (ch==' ') { continue; }

      allele[i][step] = -16;

    }
  }

  for (i=0; i<N_ALLELE; i++) {
    for (j=0; j<allele[i].size(); j++) {
      if (allele[i][j]==-1) { printf(" -1"); }  // spanning
      else if (allele[i][j] == -16) { printf(" -2"); } // loq
      else printf(" %i", allele[i][j]);
    }
    printf("\n");
  }

  */

}
