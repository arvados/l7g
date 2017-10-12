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

char rc(void) {
  int k;
  k = rand()%4;
  if (k==0) { return 'a'; }
  else if (k==1) { return 'c'; }
  else if (k==2) { return 'g'; }
  else if (k==3) { return 't'; }
  return '\0';
}

int main(int argc, char **argv) {
  int i, j, k;
  int n, z;
  double p, P=1.0/100.0;
  long int seed = -1;
  int pair, n_pair = 1;

  char *a, *b;

  n = 10000;

  if (argc>1) { n = atoi(argv[1]); }
  if (argc>2) { seed = atoi(argv[2]); }
  if (argc>3) { P = atof(argv[3]); }
  if (argc>4) { n_pair = atoi(argv[4]); }

  if (seed>=0) { srand((unsigned long int)seed); }

  a = (char *)malloc(sizeof(char)*(n+1));
  a[n] = '\0';

  for (pair=0; pair<n_pair; pair++) {

    for (i=0; i<n; i++) { a[i] = rc(); }

    printf("%s\n", a);

    for (i=0; i<n; i++) {
      p = (double)rand()/(RAND_MAX+1.0);

      if (p<P) {

        k = rand()%3;
        if (k==0) { printf("%c", rc()); }                 // sub
        else if (k==1) { continue; }                      // del
        else if (k==2) { printf("%c%c", a[i], rc()); }    // ins
      } else {
        printf("%c", a[i]);
      }
    }

    printf("\n");

  }

}
