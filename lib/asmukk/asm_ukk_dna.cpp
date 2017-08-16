#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>

#include <string>

extern "C" {
#include "asm_ukk.h"
}

extern "C" {
  int score_func(char x, char y) {
    // minimum non-zero entity (gap cost)
    //
    if ((x<0) || (y<0)) { return 2; }

    // gap
    //
    if ((x==0) || (y==0)) { return 2; }

    // 'nocall' incurs no cost to align
    //
    if ((x=='n') || (x=='N') || (y=='n') || (y=='N')) { return 0; }

    // generic mismatch
    //
    if (x!=y) { return 3; }

    // otherwise they match, no cost
    //
    return 0;
  }
}

void show_help(void) {
  printf("usage:\n");
  printf("  [-i inputfile]        Specify input file explicitely (instead of reading from stdin)\n");
  printf("  [-m mismatch_cost]    Cost of mismatched character (must be positive, default %d)\n", ASM_UKK_MISMATCH);
  printf("  [-g gap_cost]         Cost of gap (must be positive, default %d)\n", ASM_UKK_GAP);
  printf("  [-c gap_char]         Gap character (default '-')\n");
  printf("  [-S]                  Do not print aligned sequence\n");
  printf("  [-h]                  Help (this screen)\n");
}

int read_string(FILE *fp, std::string &s) {
  char ch;

  while ((ch=fgetc(fp))!=EOF) {
    if (ch=='\n') { break; }
    s += ch;
  }
  if ((ch==EOF) && (errno!=0)) { return -1; }

  return s.length();
}

int main(int argc, char **argv) {
  int k;
  char ch;
  char *input_fn = NULL;

  std::string a, b;
  char *X, *Y;
  char *a_s, *b_s;
  char gap_char = '-';

  int print_align_sequence=1;
  int mismatch_cost=ASM_UKK_MISMATCH, gap_cost=ASM_UKK_GAP;
  int score=-1;

  FILE *ifp = stdin;

  while ((ch=getopt(argc, argv, "m:g:hSi:"))!=-1) switch(ch) {
    case 'm':
      mismatch_cost = atoi(optarg);
      break;
    case 'i':
      input_fn = strdup(optarg);
      break;
    case 'g':
      gap_cost = atoi(optarg);
      break;
    case 'c':
      gap_char = optarg[0];
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

  if ((mismatch_cost<0) || (gap_cost<0)) {
    fprintf(stderr, "Mismatch cost (-m) and gap cost (-g) must both be non-zero\n");
    show_help();
    exit(1);
  }

  k = read_string(ifp, a);
  if (k<0) { perror("error reading first string"); exit(1); }
  k = read_string(ifp, b);
  if (k<0) { perror("error reading second string"); exit(1); }

  a_s = (char *)(a.c_str());
  b_s = (char *)(b.c_str());

  if (print_align_sequence) {
    score = asm_ukk_align3(&X, &Y, a_s, b_s, score_func, gap_char);
    if (score>=0) {
      printf("%d\n%s\n%s\n", score, X, Y);
    } else {
      printf("%d\n", score);
    }

    if (X) { free(X); }
    if (Y) { free(Y); }
  } else {
    //score = asm_ukk_align2(NULL, NULL, a_s, b_s, mismatch_cost, gap_cost, gap_char);
    score = asm_ukk_align3(NULL, NULL, a_s, b_s, score_func, gap_char);
    printf("%d\n", score);
  }

  return 0;
}
