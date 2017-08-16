#include <stdio.h>
#include <stdlib.h>

#include "asm_ukk.h"

#include <string>
#include <vector>

int main(int argc, char **argv) {
  FILE *fp;
  int i, k, ch;
  std::vector< std::string > ref_fns, ref_seq;
  std::string ifn, seq, tseq;
  int min_score, min_idx;

  if (argc<3) {
    printf("....\n");
    exit(-1);
  }

  for (i=1; i<(argc-1); i++) {
    ref_fns.push_back( argv[i] );
  }

  ifn = argv[argc-1];


  //DEBUG
  //
  for (i=0; i<ref_fns.size(); i++) {
    printf("%s\n", ref_fns[i].c_str());
  }
  printf("ifn: %s\n", ifn.c_str());


  for (i=0; i<ref_fns.size(); i++) {
    fp = fopen(ref_fns[i].c_str(), "r");
    if (!fp) {
      perror(ref_fns[i].c_str());
      exit(-1);
    }

    tseq.clear();
    while (!feof(fp)) {
      ch = fgetc(fp);
      if (ch==EOF)  { continue; }
      if (ch=='\n') { continue; }
      tseq += (char)ch;
    }
    ref_seq.push_back(tseq);

    fclose(fp);
  }

  fp = fopen(ifn.c_str(), "r");
  if (!fp) {
    perror(ifn.c_str());
    exit(-1);
  }
  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF)  { continue; }
    if (ch=='\n') { continue; }
    seq += (char)ch;
  }
  fclose(fp);

  for (i=0; i<ref_fns.size(); i++) {
    printf("seq(%s): %s\n",
        ref_fns[i].c_str(),
        ref_seq[i].c_str());
  }

  printf("inp: %s\n", seq.c_str());

  min_score = -1;
  min_idx = 0;
  for (i=0; i<ref_fns.size(); i++) {
    k = asm_ukk_score( (char *)ref_seq[i].c_str(), (char *)seq.c_str() );
    if ((i==0) || (k<min_score)) {
      min_score = k;
      min_idx = i;
    }
  }

  printf("min_score: %i\nmin_idx:%i\nname:%s\n",
      min_score,
      min_idx,
      ref_fns[min_idx].c_str());


}
