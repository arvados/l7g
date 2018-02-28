#include "twoBit.hpp"

// T - 00
// C - 01
// A - 10
// G - 11

twoBit_ct * twoBit_t::twoBitToCStruct() {
  twoBit_ct *twobit;

  twobit->next=NULL;
  twobit->name=NULL;
  twobit->data=NULL;
  twobit->size=0;
  twobit->nBlockCount=0;
  twobit->nStarts=NULL;
  twobit->nSizes=NULL;
  twobit->maskBlockCount=0;
  twobit->maskStarts=NULL;
  twobit->maskSizes=NULL;
  twobit->reserved=0;

  if (name.size() > 0) {
    twobit->name = strdup(name.c_str());
  }

  if (size > 0) {
    twobit->size = size;
    twobit->data = (unsigned char *)malloc(sizeof(unsigned char)*size);
    memcpy(twobit->data, &(data[0]), sizeof(unsigned char)*size);
  }

  if (nStarts.size()>0) {
    twobit->nBlockCount = nStarts.size();
    twobit->nStarts = (uint32_t *)malloc(sizeof(uint32_t)*nStarts.size());
    twobit->nSizes= (uint32_t *)malloc(sizeof(uint32_t)*nSizes.size());

    memcpy(twobit->nStarts, &(nStarts[0]), sizeof(uint32_t)*nStarts.size());
    memcpy(twobit->nSizes, &(nSizes[0]), sizeof(uint32_t)*nSizes.size());
  }

  if (maskStarts.size()>0) {
    twobit->maskBlockCount = maskStarts.size();
    twobit->maskStarts = (uint32_t *)malloc(sizeof(uint32_t)*maskStarts.size());
    twobit->maskSizes= (uint32_t *)malloc(sizeof(uint32_t)*maskSizes.size());

    memcpy(twobit->maskStarts, &(maskStarts[0]), sizeof(uint32_t)*maskStarts.size());
    memcpy(twobit->maskSizes, &(maskSizes[0]), sizeof(uint32_t)*maskSizes.size());
  }

  return twobit;
}

twoBit_ct *twoBitAlloc(void) {
  twoBit_ct *twobit;

  twobit = (twoBit_ct *)malloc(sizeof(twoBit_ct));
  twobit->next=NULL;
  twobit->name=NULL;
  twobit->data=NULL;
  twobit->size=0;
  twobit->nBlockCount=0;
  twobit->nStarts=NULL;
  twobit->nSizes=NULL;
  twobit->maskBlockCount=0;
  twobit->maskStarts=NULL;
  twobit->maskSizes=NULL;
  twobit->reserved=0;
}

void twoBitFree(twoBit_ct *twobit) {
  if (twobit) {
    if (twobit->name) { free(twobit->name); }
    if (twobit->data) { free(twobit->data); }
    if (twobit->nStarts) { free(twobit->nStarts); }
    if (twobit->nSizes) { free(twobit->nSizes); }
    if (twobit->maskStarts) { free(twobit->maskStarts); }
    if (twobit->maskSizes) { free(twobit->maskSizes); }
    free(twobit);
  }
}

int twoBit_t::twoBitFromDnaSeq(const char *seq) {
  size_t i, j, k, q, r;
  size_t n=0;
  unsigned char uc;
  uint32_t noc_s=0, noc_sz=0, mask_s=0, mask_sz=0;

  if (seq==NULL) { return -1; }

  size=0;
  name.clear();
  data.clear();
  nStarts.clear();
  nSizes.clear();
  maskStarts.clear();
  maskSizes.clear();

  while (seq[n++]);
  n--;
  if (n==0) { return 0; }
  size=n;

  data.resize((n+3)/4);

  for (i=0; i<n; i++) {
    q = i/4;
    r = i%4;

    if (r==0) { data[q] = 0; }

    uc = 0;
    if      ((seq[i]=='t') || (seq[i]=='T')) { uc = 0; }
    else if ((seq[i]=='c') || (seq[i]=='C')) { uc = 1; }
    else if ((seq[i]=='a') || (seq[i]=='A')) { uc = 2; }
    else if ((seq[i]=='g') || (seq[i]=='G')) { uc = 3; }
    data[q] |= (uc<<(2*(3-r)));

    if ((seq[i] == 'n') || (seq[i] == 'N')) {
      if (noc_sz==0) {
        noc_s = i;
        nStarts.push_back(noc_s);
      }
      noc_sz++;
    }
    else if (noc_sz>0) {
      nSizes.push_back(noc_sz);
      noc_sz=0;
    }

    if ((seq[i] == 'T') || (seq[i] == 'C') || (seq[i] == 'A') || (seq[i] == 'G') || (seq[i] == 'N')) {
      if (mask_sz==0) {
        mask_s = i;
        maskStarts.push_back(mask_s);
      }
      mask_sz++;
    }
    else if (mask_sz>0) {
      maskSizes.push_back(mask_sz);
      mask_sz=0;
    }

  }

  if (noc_sz>0) { nSizes.push_back(noc_sz); }
  if (mask_sz>0) { maskSizes.push_back(mask_sz); }

  return 0;
}

int twoBit_t::twoBitToRawDnaSeq(std::string &seq) {
  size_t i, j, k, s, n;
  size_t q, r;
  unsigned char v;

  seq.clear();

  if (size==0) { return -1; }
  n = size;

  for (i=0; i<n; i++) {
    q = i/4;
    r = i%4;
    v = (data[q] & (0x3 << (2*(3-r)))) >> (2*(3-r));
    seq += 'x';
    if      (v==0) { seq[i] = 't'; }
    else if (v==1) { seq[i] = 'c'; }
    else if (v==2) { seq[i] = 'a'; }
    else if (v==3) { seq[i] = 'g'; }
  }

  return 0;
}

int twoBit_t::twoBitToDnaSeq(std::string &seq) {
  size_t i, j, k, s, n;
  size_t q, r;
  unsigned char v;

  seq.clear();

  if (size==0) { return -1; }
  n = size;

  for (i=0; i<n; i++) {
    q = i/4;
    r = i%4;
    v = (data[q] & (0x3 << (2*(3-r)))) >> (2*(3-r));
    seq += 'x';
    if      (v==0) { seq[i] = 't'; }
    else if (v==1) { seq[i] = 'c'; }
    else if (v==2) { seq[i] = 'a'; }
    else if (v==3) { seq[i] = 'g'; }
  }

  for (i=0; i<nStarts.size(); i++) {
    for (n=0,s=nStarts[i]; n<nSizes[i]; n++) {
      seq[s+n] = 'n';
    }
  }

  for (i=0; i<maskStarts.size(); i++) {
    for (n=0,s=maskStarts[i]; n<maskSizes[i]; n++) {
      if      (seq[s+n] == 't') { seq[s+n] = 'T'; }
      else if (seq[s+n] == 'c') { seq[s+n] = 'C'; }
      else if (seq[s+n] == 'a') { seq[s+n] = 'A'; }
      else if (seq[s+n] == 'g') { seq[s+n] = 'G'; }
      else if (seq[s+n] == 'n') { seq[s+n] = 'N'; }
    }
  }

  return 0;
}

/*
void debug_print(twoBit_t *twobit) {
  int i;

  printf("next: %p\n", twobit->next);
  printf("name: %s\n", twobit->name.c_str());
  printf("data: %p\n", twobit->data);
  printf("size: %u\n", (unsigned int)twobit->size);
  printf("nBlockCount: %u\n", (unsigned int)twobit->nSizes.size());
  printf("nblock:");
  for (i=0; i<twobit->nStarts.size(); i++) {
    printf(" [%u %u]",
        (unsigned int)twobit->nStarts[i],
        (unsigned int)twobit->nSizes[i]);
  }
  printf("\n");

  printf("maskBlockCount: %u\n", (unsigned int)twobit->maskStarts.size());
  printf("maskBlocks:");
  for (i=0; i<twobit->maskStarts.size(); i++) {
    printf(" [%u %u]",
        (unsigned int)twobit->maskStarts[i],
        (unsigned int)twobit->maskSizes[i]);
  }
  printf("\n");
}
*/

void self_test(void) {
  twoBit_t twobit;
  std::string seq;

  char seqa[] = "gcatgcatgcat";
  char seqb[] = "gcanncatNNcat";
  char seqc[] = "gcatgcatgcatGCATGCATgcatgcattagctagcTaGcAcTNnnngc";
  char seqd[] = "a";

  twobit.twoBitFromDnaSeq(seqa);

  //debug_print(twobit); printf("\n");

  twobit.twoBitToDnaSeq(seq);

  if (strcmp((const char *)seq.c_str(), (const char *)seqa)==0) { printf("ok (a)\n"); }
  else { printf("ERROR on seq a\n"); }

  //printf("%s...\n", seq);

  //---

  twobit.twoBitFromDnaSeq(seqb);

  //debug_print(twobit); printf("\n");

  twobit.twoBitToDnaSeq(seq);

  if (strcmp((const char *)seq.c_str(), (const char *)seqb)==0) { printf("ok (b)\n"); }
  else { printf("ERROR on seq b\n"); }

  //printf("%s...\n", seq);

  //--

  twobit.twoBitFromDnaSeq(seqc);

  //debug_print(twobit); printf("\n");

  twobit.twoBitToDnaSeq(seq);

  if (strcmp((const char *)seq.c_str(), (const char *)seqc)==0) { printf("ok (c)\n"); }
  else { printf("ERROR on seq c\n"); }

  //printf("%s...\n", seq);

  //--

  twobit.twoBitFromDnaSeq(seqd);

  //debug_print(twobit); printf("\n");

  twobit.twoBitToDnaSeq(seq);

  if (strcmp((const char *)seq.c_str(), (const char *)seqd)==0) { printf("ok (d)\n"); }
  else { printf("ERROR on seq d\n"); }

  //printf("%s...\n", seq);

}

#ifdef TWOBIT_TEST

int main(int argc, char **argv) {
  int i, j, k;
  twoBit_t *twobit;
  unsigned char *seq;

  self_test();

  printf("ok\n");
}

#endif




