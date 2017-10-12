#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>
#include <unistd.h>
#include <errno.h>

#include <vector>
#include <string>

extern "C" {
#include "asm_ukk.h"
}

int g_mismatch = ASM_UKK_MISMATCH;
int g_gap = ASM_UKK_GAP;


extern "C" {
  int score_func(void *x, void *y, size_t sz) {
    int *ix, *iy;
    ix = (int *)x;
    iy = (int *)y;

    if ((!ix) || (!iy)) { return g_gap; }
    if ( (*ix) == (*iy) ) { return 0; }
    return g_mismatch;

  }
}

void show_help(void) {
  printf("usage:\n");
  printf("  [-i inputfile]        Specify input file explicitely (instead of reading from stdin)\n");
  printf("  [-m mismatch_cost]    Cost of mismatched character (must be positive, default %d)\n", ASM_UKK_MISMATCH);
  printf("  [-g gap_cost]         Cost of gap (must be positive, default %d)\n", ASM_UKK_GAP);
  printf("  [-c gap_int]          Gap integer (default -1)\n");
  printf("  [-S]                  Do not print aligned sequence\n");
  printf("  [-h]                  Help (this screen)\n");
}

int read_int_array(FILE *fp, std::vector< int > &x) {
  int i, v, ch;
  std::string buf;

  while ((ch=fgetc(fp))!=EOF) {
    if (ch=='\n') { break; }
    if (ch==' ') {
      if (buf.size() > 0) {
        v = atoi(buf.c_str());
        x.push_back(v);
        buf.clear();
      }
      continue;
    }

    buf += (char)ch;
  }

  if (buf.size()>0) {
    v = atoi(buf.c_str());
    x.push_back(v);
  }

  if ((ch==EOF) && (errno!=0)) { return -1; }
  return 0;
}

int main(int argc, char **argv) {
  int i, k;
  char ch;
  char *input_fn = NULL;

  int *X, *Y;
  size_t X_len, Y_len;
  std::vector< int > a, b;
  int gap_int = -1;

  int print_align_sequence=1;
  int score=-1;

  FILE *ifp = stdin;

  g_mismatch = ASM_UKK_MISMATCH;
  g_gap = ASM_UKK_GAP;

  while ((ch=getopt(argc, argv, "m:g:hSi:"))!=-1) switch(ch) {
    case 'm':
      g_mismatch = atoi(optarg);
      break;
    case 'i':
      input_fn = strdup(optarg);
      break;
    case 'g':
      g_gap = atoi(optarg);
      break;
    case 'c':
      gap_int = atoi(optarg);
      break;
    case 'S':
      print_align_sequence=0;
      break;
    case 'h':
    default:
      show_help();
      exit(0);
  }

  if ((!input_fn) && (isatty(fileno(stdin))>0)) {
    show_help();
    exit(0);
  }

  if (input_fn) {
    if ((ifp = fopen(input_fn, "r"))==NULL) {
      perror(input_fn);
      show_help();
      exit(errno);
    }
  }

  if ((g_mismatch<0) || (g_gap<0)) {
    fprintf(stderr, "Mismatch cost (-m) and gap cost (-g) must both be non-zero\n");
    show_help();
    exit(1);
  }

  k = read_int_array(ifp, a);
  if (k<0) { perror("error reading first int array"); exit(1); }
  k = read_int_array(ifp, b);
  if (k<0) { perror("error reading second int array"); exit(1); }

  if (print_align_sequence) {
    score =
      avm_ukk_align3((void **)(&X), (size_t *)(&X_len),
                     (void **)(&Y), (size_t *)(&Y_len),
                     (void *)(&(a[0])), a.size(),
                     (void *)(&(b[0])), b.size(),
                     score_func,
                     (void *)(&gap_int),
                     sizeof(int));
    if (score>=0) {
      printf("%d\n", score);
      for (i=0; i<X_len; i++) { printf(" %4d", X[i]); }
      printf("\n");
      for (i=0; i<Y_len; i++) { printf(" %4d", Y[i]); }
      printf("\n");
    } else {
      printf("%d\n", score);
    }

    if (X) { free(X); }
    if (Y) { free(Y); }
  } else {
    score =
      avm_ukk_align3(NULL, NULL,
                     NULL, NULL,
                     (void *)(&(a[0])), a.size(),
                     (void *)(&(b[0])), b.size(),
                     score_func,
                     (void *)(&gap_int),
                     sizeof(int));
    printf("%d\n", score);
  }

  return 0;
}

