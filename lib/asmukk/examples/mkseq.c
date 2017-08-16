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

#include <getopt.h>

char rc(void) {
  int k;
  k = rand()%4;
  if (k==0) { return 'a'; }
  else if (k==1) { return 'c'; }
  else if (k==2) { return 'g'; }
  else if (k==3) { return 't'; }
  return '\0';
}

void show_help(void) {
  printf("usage:\n");
  printf("  -n n              length of sequence pair\n");
  printf("  [-I ins_prob]     probability of insertion\n");
  printf("  [-D del_prob]     probability of deletion\n");
  printf("  [-U sub_prob]     probability of substitution\n");
  printf("  [-N noc_prob]     probability of nocall\n");
  printf("  [-p global_prob]  global probability (override ins., del. and sub. prob.)\n");
  printf("  [-P pair]         number of pairs to produce (default 1)\n");
  printf("  [-s seed]         seed for random number generator\n");
}

double fclamp(double x) {
  if (x<0) { return 0.0; }
  if (x>1) { return 1.0; }
  return x;
}

int main(int argc, char **argv) {
  int i, j, k;
  int n=-1, z;
  double p, P=1.0/100.0;
  long int seed = -1;
  int pair, n_pair = 1;

  double ins_prob=-1.0, del_prob=-1.0, sub_prob=-1.0, noc_prob=-1.0, glo_prob=-1.0;

  char *a, *b;

  char ch;

  double feps=1.0/1024.0;

  while ((ch=getopt(argc, argv, "n:I:D:U:s:p:N:P:h"))!=-1) switch(ch) {
    case 'n':
      n = atoi(optarg);
      break;
    case 'I':
      ins_prob = atof(optarg);
      break;
    case 'D':
      del_prob = atof(optarg);
      break;
    case 'U':
      sub_prob = atof(optarg);
      break;
    case 'N':
      noc_prob = atof(optarg);
    case 's':
      seed = atol(optarg);
      break;
    case 'p':
      glo_prob = atof(optarg);
      break;
    case 'P':
      n_pair = atoi(optarg);
      break;
    case 'h':
    default:
      show_help();
      exit(1);
      break;
  }


  if (n<=0) {
    show_help();
    exit(1);
  }

  if ((glo_prob >= 0.0) && (glo_prob<1.0)) {
    del_prob = ins_prob = sub_prob = glo_prob;
  } else {
    if ((del_prob<0) && (ins_prob<0) && (sub_prob<0)) {
      glo_prob = 1.0/(double)n;
      if ((del_prob<0) || (del_prob>1)) { del_prob=glo_prob; }
      if ((ins_prob<0) || (ins_prob>1)) { ins_prob=glo_prob; }
      if ((sub_prob<0) || (sub_prob>1)) { sub_prob=glo_prob; }
    } else {
      del_prob = fclamp(del_prob);
      ins_prob = fclamp(ins_prob);
      sub_prob = fclamp(sub_prob);
    }
  }

  if (ins_prob>=(1.0-feps)) {
    fprintf(stderr, "insertion probability cannot be 1\n");
    fflush(stderr);
    show_help();
    exit(1);
  }

  noc_prob = fclamp(noc_prob);
  //if (seed<0) { seed = 0xabecafe; }
  if (n_pair<1) { n_pair=1; }

  if (seed>=0) { srand((unsigned long int)seed); }

  a = (char *)malloc(sizeof(char)*(n+1));
  a[n] = '\0';

  for (pair=0; pair<n_pair; pair++) {

    for (i=0; i<n; i++) { a[i] = rc(); }

    printf("%s\n", a);

    for (i=0; i<n; i++) {

      // ins
      //
      p = ((double)rand()) / (RAND_MAX + 1.0);
      if (p<ins_prob) {
        printf("%c", rc());
        i--;
        continue;
      }

      // del
      //
      p = ((double)rand()) / (RAND_MAX + 1.0);
      if (p<del_prob) { continue; }

      // sub
      //
      p = ((double)rand()) / (RAND_MAX + 1.0);
      if (p<sub_prob) {
        printf("%c", rc());
        continue;
      }

      // nocall
      //
      p = ((double)rand()) / (RAND_MAX + 1.0);
      if (p<noc_prob) {
        printf("n");
        continue;
      }

      // default
      //
      printf("%c", a[i]);
    }

    printf("\n");

  }

}
