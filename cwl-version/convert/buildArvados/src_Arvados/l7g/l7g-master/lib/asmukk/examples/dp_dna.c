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

#include <vector>
#include <string>

int g_debug=0;

int default_score(char x, char y) {
  if ((x<=0) || (y<=0)) { return 2; }
  if ((x=='n') || (x=='N') || (y=='n') || (y=='N')) { return 0; }
  if (x!=y) { return 3; }
  return 0;
}

int default_gap(char x) {
  return 2;
}

void dp_D_print(int *D, const char *a, const char *b, int n_c, int m_r) {
  int i, j, k;

  printf("   ");
  for (i=0; i<n_c; i++) {
    if (i==0) { printf("  -"); }
    else { printf(" %2c", a[i-1]); }
  }
  printf("\n");

  for (j=0; j<m_r; j++) {
    if (j==0) { printf("  -"); }
    else { printf(" %2c", b[j-1]); }
    for (i=0; i<n_c; i++) {
      printf(" %2i", D[j*n_c + i]);
    }
    printf("\n");
  }
  printf("\n");
}

void dp_D_print3(int *D, const char *a, const char *b, int n_c, int m_r) {
  int i, j, k;

  printf("    ");
  for (i=0; i<n_c; i++) {
    if (i==0) { printf("   -"); }
    else { printf("  %2c", a[i-1]); }
  }
  printf("\n");

  for (j=0; j<m_r; j++) {
    if (j==0) { printf("   -"); }
    else { printf("  %2c", b[j-1]); }
    for (i=0; i<n_c; i++) {
      printf(" %3i", D[j*n_c + i]);
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
int dp_score(const char *a, const char *b, int (*score)(char,char), int (*gap_a)(char), int (*gap_b)(char)) {
  int i, j, k;
  int n_c, m_r;
  int *D, d;

  n_c = strlen(a)+1;
  m_r = strlen(b)+1;

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

  //DEBUG
  //dp_D_print3(D, a, b, n_c, m_r);


  d = D[ (m_r-1)*n_c + n_c-1 ];
  free(D);
  return d;
}

int dp_simp(const char *a, const char *b) {
  return dp_score(a, b, default_score, default_gap, default_gap);
}

int dp_align(char **X, char **Y, const char *a, const char *b, int (*score)(char,char), int (*gap_a)(char), int (*gap_b)(char)) {
  int i, j, k;
  int n_c, m_r;
  int *D, d;
  int align_len=0;
  int cur_r, cur_c, cur_val, cur_pos;
  int dr, dc;

  n_c = strlen(a)+1;
  m_r = strlen(b)+1;

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


  *X = (char *)malloc(sizeof(char)*(align_len+1));
  *Y = (char *)malloc(sizeof(char)*(align_len+1));

  // Calculate alignment
  //
  cur_r = m_r-1;
  cur_c = n_c-1;
  cur_pos = align_len;
  (*X)[cur_pos] = '\0';
  (*Y)[cur_pos] = '\0';
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
      (*X)[cur_pos] = '-';
      (*Y)[cur_pos] = b[cur_r];
    }
    else {
      (*X)[cur_pos] = a[cur_c];
      (*Y)[cur_pos] = '-';
    }

  }

  if (g_debug) { dp_D_print(D, a, b, n_c, m_r); }

  d = D[ (m_r-1)*n_c + n_c-1 ];
  free(D);
  return d;
}

int dp_align_simple(char **X, char **Y, const char *a, const char *b) {
  dp_align(X,Y,a,b,default_score,default_gap, default_gap);
}

int main(int argc, char **argv) {
  char ch;
  int sc, sc_align;
  char *X, *Y;
  int line_counter=0;
  std::string a, b;

  while ((ch=fgetc(stdin))!=EOF) {
    if (ch=='\n') {
      line_counter++;
      if (line_counter==2) { break; }
      continue;
    }

    if (line_counter==0) {
      a += (char)ch;
    } else {
      b += (char)ch;
    }
  }


  sc       = dp_simp(a.c_str(), b.c_str());
  sc_align = dp_align_simple(&X, &Y, a.c_str(), b.c_str());

  //printf("%d (%d)\n", sc, sc_align);
  printf("%d\n", sc);
  printf("%s\n%s\n", X, Y);
}
