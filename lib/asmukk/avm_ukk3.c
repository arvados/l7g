/*
    Copyright (C) 2017 Curoverse, Inc.

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

#include "asm_ukk.h"

static int g_verbose=0;
static int g_debug=0;

// void data compare
//
// -1 a  < b (lex ordering)
// +1 a  > b (lex ordering)
//  0 a == b
//
static int vd_cmp(void *a, void *b, size_t n) {
  size_t i;
  char *c_a, *c_b;
  c_a = (char *)a;
  c_b = (char *)b;
  for (i=0; i<n; i++) {
    if (c_a[i]!=c_b[i]) {
      if (c_a[i]<c_b[i]) { return -1; }
      return 1;
    }
  }
  return 0;
}

static int default_vd_score_func(void *a, void *b, size_t n) {
  if ((!a) || (!b)) { return 2; }
  if (vd_cmp(a,b,n)==0) { return 0; }
  return 3;
}

static void byte_swap(void *a, void *b, size_t sz) {
  size_t i;
  unsigned char ch, *c_a, *c_b;
  c_a = a;
  c_b = b;
  for (i=0; i<sz; i++) {
    ch = c_a[i]; c_a[i] = c_b[i]; c_b[i] = ch;
  }
}

int align_v_W3(void **X, size_t *X_len,
               void **Y, size_t *Y_len,
               void *a, size_t a_len,
               void *b, size_t b_len,
               int *W, int m_r, int n_c,
               int w_len,
               int (*score_func)(void *, void *, size_t),
               void *gap_ele,
               size_t sz,
               int row_pref) {
  int i;
  int dr, dc;
  int r, c, w;
  int pos00, pos01, pos10, pos11;
  int w_offset;
  int mm;
  size_t xy_pos=0;
  char ch;

  char *tx, *ty;

  i = ((n_c>m_r)?n_c:m_r);

  *X = NULL;
  *Y = NULL;

  tx = (void *)malloc(sz*2*i);
  ty = (void *)malloc(sz*2*i);

  w_offset = w_len/2;

  r = m_r-1;
  c = n_c-1;
  while ((r>0) || (c>0)) {
    dr = 0;
    dc = 0;

    w = c - (r-w_offset);
    pos11 = r*w_len + w;

    // if `row_pref` is set the preference is:
    //   * straight alignment
    //   * row
    //   * col
    //
    // otherwise if `row_pref` is not set, the
    // preference is:
    //   * straight alignment
    //   * col
    //   * row
    //
    // Conditionals lower down overwrite decisions above, so the
    // lower the conditional, the higher the precedence.
    //

    if (row_pref) {

      if (c>0) {
        w = (c-1) - (r-w_offset);
        if ((w>=0) && (w<w_len)) {
          pos10 = r*w_len + w;
          if ((W[pos10]+score_func(a + sz*(c-1), NULL, sz)) == W[pos11]) { dr=0; dc=-1; }
        }
      }

      if (r>0) {
        w = c - ((r-1)-w_offset);
        if ((w>=0) && (w<w_len)) {
          pos01 = (r-1)*w_len + w;
          if ((W[pos01]+score_func(NULL, b + sz*(r-1), sz)) == W[pos11]) { dr=-1; dc=0; }
        }
      }

    } else {

      if (r>0) {
        w = c - ((r-1)-w_offset);
        if ((w>=0) && (w<w_len)) {
          pos01 = (r-1)*w_len + w;
          if ((W[pos01]+score_func(NULL, b + sz*(r-1), sz)) == W[pos11]) { dr=-1; dc=0; }
        }
      }

      if (c>0) {
        w = (c-1) - (r-w_offset);
        if ((w>=0) && (w<w_len)) {
          pos10 = r*w_len + w;
          if ((W[pos10]+score_func(a + sz*(c-1), NULL, sz)) == W[pos11]) { dr=0; dc=-1; }
        }
      }

    }

    if ((r>0) && (c>0)) {
      w = (c-1) - ((r-1)-w_offset);
      if ((w>=0) && (w<w_len)) {
        pos00 = (r-1)*w_len + w;
        mm = score_func(a + sz*(c-1), b + sz*(r-1), sz);
        if ((W[pos00]+mm) == W[pos11]) { dr=-1; dc=-1; }
      }
    }

    if ((dr==-1) && (dc==-1)) {

      memcpy(tx + sz*xy_pos, a + sz*(c-1), sz);
      memcpy(ty + sz*xy_pos, b + sz*(r-1), sz);

    } else if ((dr==-1) && (dc==0)) {

      memcpy(tx + sz*xy_pos, gap_ele, sz);
      memcpy(ty + sz*xy_pos, b + sz*(r-1), sz);

    } else if ((dr==0) && (dc==-1)) {

      memcpy(tx + sz*xy_pos, a + sz*(c-1), sz);
      memcpy(ty + sz*xy_pos, gap_ele, sz);

    } else {

      free(tx);
      free(ty);

      return -1;
    }

    xy_pos++;
    r+=dr;
    c+=dc;
  }

  *X_len = xy_pos;
  *Y_len = xy_pos;

  for (i=0; i<(xy_pos/2); i++) {
    byte_swap(tx + sz*i, tx + sz*(xy_pos-i-1), sz);
    byte_swap(ty + sz*i, ty + sz*(xy_pos-i-1), sz);
  }

  *X = tx;
  *Y = ty;

  return 0;
}

int avm_ukk_score3(void *a, int a_len,
                   void *b, int b_len,
                   int (*score_func)(void *, void *, size_t),
                   size_t sz) {
  int threshold = (1<<2);
  int it, max_it=((32-2-1));
  int sc = -2;

  for (it=0; (it<max_it) && (sc<0); it++) {
    sc = vd_align_ukk3(NULL, NULL, NULL, NULL, a, a_len, b, b_len, threshold, score_func, NULL, sz);
    threshold*=2;
  }

  return sc;
}

// Run Ukkonnen's approximate string alignment on `a` and `b`
//   storing result in `X` and `Y` using
// `score_func` as the scoring function and `gap_ele` as the gap data.
// vd_align_ukk3 is called with a threshold that is doubled after every failed
//   alignment.
// `X` and `Y` are arrays of size `sz`
//
int avm_ukk_align3(void **X, size_t *X_len,
                   void **Y, size_t *Y_len,
                   void *a, size_t a_len,
                   void *b, size_t b_len,
                   int (*score_func)(void *, void *, size_t),
                   void *gap_ele,
                   size_t sz) {
  int threshold = (1<<2);
  int it, max_it=((32-2-1));
  int sc = -2;

  if ((X!=NULL) && (Y!=NULL)) {
    *X = NULL;
    *Y = NULL;
  }

  for (it=0; (it<max_it) && (sc<0); it++) {

    if (g_verbose) { printf("# threshold %d\n", threshold); }

    sc = vd_align_ukk3(X, X_len, Y, Y_len, a, a_len, b, b_len, threshold, score_func, gap_ele, sz);
    threshold*=2;

    if (sc<0) {
      if ((X!=NULL) && (Y!=NULL)) {
        if (*X) free(*X);
        if (*Y) free(*Y);
        *X = NULL;
        *Y = NULL;
      }
    }
  }

  return sc;
}

// Run Ukkonnen's approximate string alignment on `a` and `b` up until
//   threshold T has been reached, storing result in `X` and `Y` using
// ` score_func` as the scoring function and `gap_char` as the gap character.
// -1 is returned if theshold `T` was reached.
//
int vd_align_ukk3(void **X, size_t *X_len,
                  void **Y, size_t *Y_len,
                  void *a_orig, size_t a_len,
                  void *b_orig, size_t b_len,
                  int32_t T,
                  int (*score_func)(void *, void *, size_t),
                  void *gap_ele, size_t sz) {
  int ret;
  ssize_t r,c, n_c, m_r, len_ovf;
  int *W;
  ssize_t w, w_offset, w_len;
  ssize_t p, del, m;
  int create_align_seq = 0;

  int i, j;
  void *a, *b, *TXY=NULL;
  int seq_swap=0;

  int gap_cost;

  n_c = (int)a_len+1;
  m_r = (int)b_len+1;

  a = a_orig;
  b = b_orig;

  if (m_r > n_c) {
    a = b_orig;
    b = a_orig;
    i = n_c; n_c=m_r; m_r=i;
    seq_swap=1;
  }

  // Find minimum non-zero score for
  // window band space allocation.
  //
  del = score_func(NULL, NULL, sz);
  if (del<=0) { return -1; }

  gap_cost = score_func(NULL, NULL, sz);

  if (X && Y) { create_align_seq = 1; }

  if (create_align_seq) {
    *X = NULL;
    *Y = NULL;
  }

  // t/del < |n-m| -> reject
  //
  len_ovf = ((n_c>m_r) ? (n_c-m_r) : (m_r-n_c));
  if ((T/del) < len_ovf) {
    if (create_align_seq) {
      if (!(*X)) free(*X);
      if (!(*Y)) free(*Y);
    }
    return -1;
  }

  p = (T/del) - len_ovf;
  p /= 2;

  w_offset = ((n_c>m_r) ? (n_c-m_r+p) : p);
  w_len = 2*w_offset+1;

  // our window isn't big enough to hold calculated values
  //
  w = (n_c-1) - ((m_r-1)-w_offset);
  if ((w<0) || (w>=w_len)) {
    if (create_align_seq) {
      if (!(*X)) free(*X);
      if (!(*Y)) free(*Y);
    }
    return -1;
  }

  W = (int *)malloc(sizeof(int)*m_r*w_len);
  if (!W) {
    fprintf(stderr, "could not allocate matrix\n");
    fflush(stderr);
    return -1;
  }

  for (w=0; w<w_len; w++) {
    c = w-w_offset;

    if (w<w_offset) { W[w] = -1; }
    else { W[w] = gap_cost*(w-w_offset); }

  }

  for (r=1; r<m_r; r++) {

    // For conceptual simplicity, enumerate columns
    //
    for (c=(r-w_offset); c<=(r+w_offset); c++) {


      // Window position
      //
      w = c - (r-w_offset);
      int w_rmm = c - ((r-1)-w_offset);

      if (c<0) { W[r*w_len + w] = -1; }

      else if (c==0) { W[r*w_len + w] = W[(r-1)*w_len + w_rmm] + score_func(NULL, b + sz*(r-1), sz); }

      else if (c>=n_c) { W[r*w_len + w] = -1; }
      else {

        // diagonal value
        //
        m = W[(r-1)*w_len + w] + score_func(a + sz*(c-1),b + sz*(r-1), sz) ;


        // left to right transition
        //
        if ((w>0) && ((W[r*w_len+w-1] + score_func(NULL, b + sz*(r-1), sz)) < m)) { m = W[r*w_len+w-1] + score_func(NULL, b + sz*(r-1), sz); }


        // top to bottom transition
        //
        if ((w+1)!=w_len) {
          if ((W[(r-1)*w_len+w+1] + score_func(a + sz*(c-1), NULL, sz)) < m) { m = W[(r-1)*w_len+w+1] + score_func(a + sz*(c-1), NULL, sz); }
        }

        W[r*w_len+w] = m;
      }

    }
  }

  w = (n_c-1) - ((m_r-1)-w_offset);
  m = W[(m_r-1)*w_len + w];

  if (m>T) {
    if (create_align_seq) {
      if (!(*X)) free(*X);
      if (!(*Y)) free(*Y);
    }
    return -1;
  }


  if (create_align_seq) {
    ret = align_v_W3(X, X_len, Y, Y_len, a, a_len, b, b_len, W, m_r, n_c, w_len, score_func, gap_ele, sz, seq_swap);
    if (ret<0) { return ret; }
  }

  free(W);

  if (create_align_seq && seq_swap) {
    TXY = *X;
    *X   = *Y;
    *Y   = TXY;
  }

  return m;
}
