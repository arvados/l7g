/*
    Copyright (C) 2015 Curoverse, Inc.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>

#include <getopt.h>

#include <vector>
#include <string>

int g_debug=0;

int g_mismatch_cost = 3;
int g_gap_cost = 2;

int default_score(int x, int y) {
  if (x==y) { return 0; }
  return g_mismatch_cost;
}

int default_gap(int x) {
  return g_gap_cost;
}

void dp_D_print(int *D, int *a, size_t a_len, int *b, size_t b_len, int n_c, int m_r) {
  int i, j, k;

  printf("   ");
  for (i=0; i<n_c; i++) {
    if (i==0) { printf("  -"); }
    else { printf(" %4i", a[i-1]); }
  }
  printf("\n");

  for (j=0; j<m_r; j++) {
    if (j==0) { printf("  -"); }
    else { printf(" %4i", b[j-1]); }
    for (i=0; i<n_c; i++) {
      printf(" %4i", D[j*n_c + i]);
    }
    printf("\n");
  }
  printf("\n");
}

int min3(int x, int y, int z) {
  int m;
  m = x;
  if (y<m) { m = y; }
  if (z<m) { m = z; }
  return m;
}

// a on columns
// b on rows
//
int dp_score(int *a, size_t a_len, int *b, size_t b_len, int (*score)(int,int), int (*gap_a)(int), int (*gap_b)(int)) {
  int i, j, k;
  int n_c, m_r;
  int *D, d;

  //n_c = strlen(a)+1;
  //m_r = strlen(b)+1;

  n_c = a_len+1;
  m_r = b_len+1;

  D = (int *)malloc(sizeof(int)*n_c*m_r);

  D[0] = 0;
  for (i=1; i<n_c; i++) {
    D[i] = D[i-1] + gap_a(a[i-1]);
  }

  for (i=1; i<m_r; i++) {
    D[i*n_c] = D[(i-1)*n_c] + gap_b(b[i-1]);
  }

  for (i=1; i<n_c; i++) {
    for (j=1; j<m_r; j++) {
      D[j*n_c + i] = min3( D[(j-1)*n_c + (i-1)] + score(a[i-1], b[j-1]),
                           D[(j-1)*n_c + i] + gap_a(a[i]),
                           D[j*n_c + (i-1)] + gap_b(b[j]) );
    }
  }

  d = D[ (m_r-1)*n_c + n_c-1 ];
  free(D);
  return d;
}

int dp_simp(int *a, size_t a_len, int *b, size_t b_len) {
  return dp_score(a, a_len, b, b_len, default_score, default_gap, default_gap);
}

int dp_align(int **X, size_t *X_len,
             int **Y, size_t *Y_len,
             int *a, size_t a_len,
             int *b, size_t b_len,
             int (*score)(int,int),
             int (*gap_a)(int),
             int (*gap_b)(int)) {
  int i, j, k;
  int n_c, m_r;
  int *D, d;
  int align_len=0;
  int cur_r, cur_c, cur_val, cur_pos;
  int dr, dc;

  //n_c = strlen(a)+1;
  //m_r = strlen(b)+1;

  n_c = a_len+1;
  m_r = b_len+1;

  D = (int *)malloc(sizeof(int)*n_c*m_r);

  D[0] = 0;
  for (i=1; i<n_c; i++) {
    D[i] = D[i-1] + gap_a(a[i-1]);
  }

  for (i=1; i<m_r; i++) {
    D[i*n_c] = D[(i-1)*n_c] + gap_b(b[i-1]);
  }

  for (j=1; j<m_r; j++) {
    for (i=1; i<n_c; i++) {
      D[j*n_c + i] = min3( D[(j-1)*n_c + (i-1)] + score(a[i-1], b[j-1]),
                           D[(j-1)*n_c + i] + gap_a(a[i]),
                           D[j*n_c + (i-1)] + gap_b(b[j]) );
    }
  }

  // calculate length
  //
  align_len = 0;
  cur_r = m_r-1;
  cur_c = n_c-1;
  while ((cur_r>0) || (cur_c>0)) {
    dr = 0;
    dc = 0;
    cur_val = D[cur_r*n_c + cur_c];
    if ((cur_r>0) &&  (cur_val == (D[(cur_r-1)*n_c + cur_c] + gap_a(b[cur_r-1])))) {
      dr=-1; dc = 0;
    }
    if ((cur_c>0) &&  (cur_val == (D[cur_r*n_c + (cur_c-1)] + gap_b(a[cur_c-1])))) {
      dr=0; dc = -1;
    }
    if ((cur_r>0) && (cur_c>0) &&  (cur_val == (D[(cur_r-1)*n_c + (cur_c-1)] + score(a[cur_c-1], b[cur_r-1])))) {
      dr=-1; dc=-1;
    }

    cur_r+=dr;
    cur_c+=dc;
    align_len++;
  }


  *X = (int *)malloc(sizeof(int)*(align_len+1));
  *Y = (int *)malloc(sizeof(int)*(align_len+1));

  // Calculate alignment
  //
  cur_r = m_r-1;
  cur_c = n_c-1;
  cur_pos = align_len;

  //(*X)[cur_pos] = '\0';
  //(*Y)[cur_pos] = '\0';
  *X_len = cur_pos;
  *Y_len = cur_pos;

  while ((cur_r>0) || (cur_c>0)) {
    dr = 0;
    dc = 0;
    cur_val = D[cur_r*n_c + cur_c];
    if ((cur_r>0) &&  (cur_val == (D[(cur_r-1)*n_c + cur_c] + gap_b(b[cur_r-1])))) {
      dr=-1; dc = 0;
    }
    if ((cur_c>0) &&  (cur_val == (D[cur_r*n_c + (cur_c-1)] + gap_a(a[cur_c-1])))) {
      dr=0; dc = -1;
    }
    if ((cur_r>0) && (cur_c>0) &&  (cur_val == (D[(cur_r-1)*n_c + (cur_c-1)] + score(a[cur_c-1], b[cur_r-1])))) {
      dr=-1; dc=-1;
    }

    cur_r+=dr;
    cur_c+=dc;
    cur_pos--;

    if ((dr==-1) && (dc==-1)) {
      (*X)[cur_pos] = a[cur_c];
      (*Y)[cur_pos] = b[cur_r];
    }
    else if (dr==-1) {
      (*X)[cur_pos] = -1;
      (*Y)[cur_pos] = b[cur_r];
    }
    else {
      (*X)[cur_pos] = a[cur_c];
      (*Y)[cur_pos] = -1;
    }

  }

  if (g_debug) { dp_D_print(D, a, a_len, b, b_len, n_c, m_r); }

  d = D[ (m_r-1)*n_c + n_c-1 ];
  free(D);
  return d;
}

int dp_align_simple(int **X, size_t *X_len, int **Y, size_t *Y_len, int *a, size_t a_len, int *b, size_t b_len) {
  dp_align(X, X_len,
           Y, Y_len,
           a, a_len,
           b, b_len,
           default_score,
           default_gap,
           default_gap);
}

void show_help(void) {
  printf("usage:\n");
  printf("\n");
  printf("    dpi [-m mismatch] [-g gap] [-h] < seqfn\n");
  printf("\n");
  printf("  [-h]            show help (this scree)\n");
  printf("  [-m mismatch]   mismatch cost (default 3)\n");
  printf("  [-g gap]        gap cost (default 2)\n");
  printf("\n");
}

int main(int argc, char **argv) {
  char ch;
  int sc, sc_align;
  int i;

  std::vector< int > a, b;
  int *X, *Y;
  size_t X_len, Y_len;

  int line_counter=0;
  std::vector< int > u, v;
  std::string buf;

  while ((ch=getopt(argc, argv, "hm:g:"))!=-1) switch (ch) {
    case 'm':
      g_mismatch_cost = atoi(optarg); break;
    case 'g':
      g_gap_cost = atoi(optarg); break;
    default:
    case 'h':
      show_help(); exit(0); break;
  }

  while ((ch=fgetc(stdin))!=EOF) {
    if (ch=='\n') {
      if (buf.size()>0) {
        if (line_counter==0) { a.push_back(atoi(buf.c_str())); }
        else { b.push_back(atoi(buf.c_str())); }
      }
      buf.clear();

      line_counter++;
      if (line_counter==2) { break; }
      continue;
    }

    if (ch==' ') {
      if (buf.size()>0) {
        if (line_counter==0) { a.push_back(atoi(buf.c_str())); }
        else { b.push_back(atoi(buf.c_str())); }
      }
      buf.clear();
      continue;
    }
    buf += (char)ch;
  }
  if (buf.size()>0) {
    if (line_counter==0) { a.push_back(atoi(buf.c_str())); }
    else { b.push_back(atoi(buf.c_str())); }
  }

  sc       = dp_simp(&(a[0]), a.size(), &(b[0]), b.size());
  sc_align = dp_align_simple(&X, &X_len, &Y, &Y_len, &(a[0]), a.size(), &(b[0]), b.size());

  //printf("%d (%d)\n", sc, sc_align);
  printf("%d\n", sc);
  for (i=0; i<X_len; i++) { printf(" %4i", X[i]); }
  printf("\n");
  for (i=0; i<Y_len; i++) { printf(" %4i", Y[i]); }
  printf("\n");

}
