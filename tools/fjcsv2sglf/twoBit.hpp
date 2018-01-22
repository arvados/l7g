#ifndef TWOBIT_H
#define TWOBIT_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <string.h>
#include <strings.h>

#include <string>
#include <vector>

//#define TWOBIT_TEST


// T - 00
// C - 01
// A - 10
// G - 11

/* Two bit representation of DNA. */
typedef struct twoBit_ctype {
  struct twoBit_ctype *next;  /* Next sequence in list */
  char *name;     /* Name of sequence. */
  unsigned char *data;    /* DNA at two bits per base. */
  uint32_t size;    /* Size of this sequence. */
  uint32_t nBlockCount;   /* Count of blocks of Ns. */
  uint32_t *nStarts;    /* Starts of blocks of Ns. */
  uint32_t *nSizes;   /* Sizes of blocks of Ns. */
  uint32_t maskBlockCount;  /* Count of masked blocks. */
  uint32_t *maskStarts;   /* Starts of masked regions. */
  uint32_t *maskSizes;    /* Sizes of masked regions. */
  uint32_t reserved;    /* Reserved for future expansion. */
} twoBit_ct;

/* Two bit representation of DNA. */
typedef struct twoBit_type {
  struct twoBit_type *next;  /* Next sequence in list */
  std::string name;     /* Name of sequence. */
  uint32_t size;
  std::vector< unsigned char > data;
  std::vector< uint32_t > nStarts;
  std::vector< uint32_t > nSizes;
  std::vector< uint32_t > maskStarts;
  std::vector< uint32_t > maskSizes;
  uint32_t reserved;    /* Reserved for future expansion. */

  twoBit_type() : next(NULL) { }
  twoBit_ct *twoBitToCStruct();
  int twoBitFromDnaSeq(const char *seq);
  int twoBitToDnaSeq(std::string &seq);
  int twoBitToRawDnaSeq(std::string &seq);

} twoBit_t;


//twoBit_t *twoBitFromDnaSeq(twoBit_t *twobit, const char *seq);
unsigned char *twoBitToDnaSeq_c(twoBit_t *twobit, unsigned char *seq);
unsigned char *twoBitToDnaSeqStr(twoBit_t *twobit, std::string &seq);

twoBit_ct *twoBitAlloc(void);
void twoBitFree(twoBit_ct *twobit);

#endif
