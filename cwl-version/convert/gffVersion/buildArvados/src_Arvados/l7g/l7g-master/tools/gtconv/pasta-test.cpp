#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

#include <string>

#include "pasta.hpp"

int main(int argc, char **argv) {
  int i, j, k;
  int ch;

  char b[] = "acgtn";
  std::string ref, inp, expected;
  std::string pasta_str;

  std::string ref_check, inp_check;

  for (i=0; i<5; i++) {
    ref += b[i];
    inp += b[i];
    expected += b[i];
  }

  for (i=0; i<5; i++) {
    ref += b[i];
    inp += b[(i+1)%5];
  }
  expected += ".....";

  for (i=0; i<5; i++) {
    ref += b[i];
    inp += b[(i+2)%5];
  }
  expected += ".....";

  for (i=0; i<5; i++) {
    ref += b[i];
    inp += b[(i+3)%5];
  }
  expected += ".....";

  for (i=0; i<5; i++) {
    ref += b[i];
    inp += b[(i+4)%5];
  }
  expected += ".....";

  for (i=0; i<5; i++) {
    ref += b[i];
    inp += " ";
  }
  expected += ".....";

  for (i=0; i<5; i++) {
    ref += " ";
    inp += b[i];
  }
  expected += ".....";

  for (i=0; i<ref.size(); i++) {
    if (ref[i]==' ') {
      ch = pasta_convert(0, inp[i]);
    }
    else if (inp[i] == ' ') {
      ch = pasta_convert(ref[i], 0);
    }
    else {
      ch = pasta_convert(ref[i], inp[i]);
    }
    printf(" %c %c -> %c (%c)\n",
        ref[i], inp[i], expected[i], (char)ch);

    pasta_str += (char)ch;
  }

  printf("pasta seq: %s\n", pasta_str.c_str());

  for (i=0; i<pasta_str.size(); i++) {
    ch = pasta2ref(pasta_str[i]);
    if (ch==0) { ref_check += ' '; }
    else { ref_check += (char)ch; }

    ch = pasta2seq(pasta_str[i]);
    if (ch==0) { inp_check += ' '; }
    else { inp_check += (char)ch; }
  }

  if (ref != ref_check) {
    fprintf(stderr, "reference sequences don't match!\n");
    exit(-1);
  }

  if (inp != inp_check) {
    fprintf(stderr, "sequences don't match!\n");
    exit(-2);
  }

  printf("ref: %s\nchk: %s\n", ref.c_str(), ref_check.c_str());
  printf("seq: %s\nchk: %s\n", inp.c_str(), inp_check.c_str());

  exit(0);

}
