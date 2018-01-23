#include <stdio.h>
#include <stdlib.h>

#include <string>
#include <vector>
#include <map>
#include <algorithm>

#include <openssl/md5.h>

#include "twoBit.hpp"

//#define DEBUG_RUN

// Container element for the twobit sequence and frequency
//
typedef struct sglf2bit_type {
  twoBit_t *twobit;
  int freq;
  sglf2bit_type() : twobit(NULL), freq(0) { }
} sglf2bit_t;

typedef struct opt_type {
  int dedup_fj;
} opt_t;

bool twobit_cmp(twoBit_t *a, twoBit_t *b) {
  if (a->name != b->name) {
    return a->name < b->name;
  }
  if (a->size != b->size) { return a->size < b->size; }
  return false;
}

// Compare on (normalized) tile ID name and then
// on sequence size.
//
bool sglf2bit_cmp(sglf2bit_t *a, sglf2bit_t *b) {
  if (a->twobit->name != b->twobit->name) {
    return a->twobit->name < b->twobit->name;
  }
  if (a->twobit->size != b->twobit->size) { return a->twobit->size < b->twobit->size; }
  return false;
}

void print_sglf2bit(std::vector<sglf2bit_t *> &sglf2bit) {
  int i, j, k;
  std::string seq;

  for (i=0; i<sglf2bit.size(); i++) {
    sglf2bit[i]->twobit->twoBitToDnaSeq(seq);
    printf("%s,%s\n", sglf2bit[i]->twobit->name.c_str(), seq.c_str());
  }
}

// For ease of processing, normalize the tile ID
//
int normalize_name(twoBit_t *twobit) {
  int i, j, state=0, i_n;
  size_t n;
  std::string new_name;
  int prev_pos=0;

  if (twobit->name.size()==0) {

    // place holder name, not necessarily a real tileid or a
    // tileid we are using.
    //
    twobit->name = "000f.00.000f.000f+000f";
    return 0;
  }
  n = twobit->name.size();

  for (i=0; i<n; i++) {
    if ((twobit->name[i]=='.') || (twobit->name[i]=='+')) {
      if ((state==0) || (state==2) || (state==3)) {

        // zero out variant
        //
        if (state==3) {
          for (j=0; j<4; j++) { new_name += '0'; }
          for (; prev_pos<i; prev_pos++) ;
        }

        // otherwise pad at '0' at the beginning (4 length)
        //
        else {
          for (j=0; j<(4-(i-prev_pos)); j++) { new_name += '0'; }
          for (; prev_pos<i; prev_pos++) { new_name += twobit->name[prev_pos]; }
        }

        // move past separator
        //
        prev_pos++;

        // add separator to new_name
        //
        if (state<3) { new_name += '.'; }
        else if (state<4) { new_name += '+'; }
        else { return -1; }
      }

      // lib version
      //
      else if (state==1) {
        for (j=0; j<(2-(i-prev_pos)); j++) { new_name += '0'; }
        for (; prev_pos<i; prev_pos++) { new_name += twobit->name[prev_pos]; }
        prev_pos++;
        if (state<3) { new_name += '.'; }
        else { return -2; }
      }
      else { return -3; }
      state++;

      continue;
    }
  }

  i_n = (int)n;

  if (state!=4) { return -4; }
  for (j=0; j<(4-(i_n-prev_pos)); j++) { new_name += '0'; }
  for (; prev_pos<i_n; prev_pos++) { new_name += twobit->name[prev_pos]; }

  twobit->name = new_name.c_str();

  return 1;
}

// Read in FastJ stream and store in sglf2bit vector.
// Sort on tile ID and sequence length for later processing.
//
int read_csvseq_to_twobit(FILE *ifp, std::vector<sglf2bit_t *> &sglf2bit, opt_t &opt) {
  int i, ch;
  sglf2bit_t *s2b=NULL;
  twoBit_t *twobit=NULL;
  int state=0, tile_freq;
  std::string name, m5, seq;

  std::map< std::string, int > fj_m5_map, fj_m5_idx_map;
  std::map< std::string, int >::iterator srch;

#ifdef DEBUG_RUN
  uint64_t line_no=0, line_print_n = 1000;
#endif

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if ((ch=='\n') || (ch==EOF)) {

#ifdef DEBUG_RUN
      line_no++;
      if ((line_no%line_print_n)==0) {
        fprintf(stderr, "# line_no %llu\n", (unsigned long long int)line_no);
        fflush(stderr);
      }
#endif

      if (seq.size()==0) {
        name.clear();
        m5.clear();
        seq.clear();
        state=0;
        continue;
      }

      // only insert unique fj tile if option set
      //
      if (opt.dedup_fj) {

        tile_freq = 0;
        srch = fj_m5_map.find(m5);
        if (srch != fj_m5_map.end()) {

          tile_freq = srch->second;
        }
        else { }

        tile_freq++;
        fj_m5_map[m5] = tile_freq;

        if (tile_freq == 1) {
          twobit = new twoBit_t;
          twobit->twoBitFromDnaSeq((const char *)seq.c_str());

          fj_m5_idx_map[m5] = (int)sglf2bit.size();

          if (name.size()>0) {
            twobit->name = name;
            normalize_name(twobit);
          }

          s2b = new sglf2bit_t;
          s2b->freq = tile_freq;
          s2b->twobit = twobit;
          sglf2bit.push_back(s2b);

        }
        else {

          srch = fj_m5_idx_map.find(m5);
          if (srch != fj_m5_idx_map.end()) {
            sglf2bit[ srch->second ]->freq = tile_freq;
          }

        }

        name.clear(); m5.clear(); seq.clear();
        state=0;
        continue;
      }

      // default case
      //
      twobit = new twoBit_t;
      twobit->twoBitFromDnaSeq((const char *)seq.c_str());


      if (name.size()>0) {
        twobit->name = name;
        normalize_name(twobit);
      }

      s2b = new sglf2bit_t;
      s2b->freq = 1;
      s2b->twobit = twobit;
      sglf2bit.push_back(s2b);

      name.clear();
      m5.clear();
      seq.clear();
      state=0;
      continue;
    }

    if (ch==',') { state++; continue; }

    if      (state==0) { name += (char)ch; }
    else if (state==1) { m5 += (char)ch; }
    else if (state==2) { seq += (char)ch; }
    else { return -1; }

  }
  twobit=NULL;

  std::sort(sglf2bit.begin(), sglf2bit.end(), sglf2bit_cmp);

  return 0;
}

// Bubble sort on first vector, freq, and mirror swaps
// on vector `val`
//
void bub_mirror(std::vector<int> &freq, std::vector<char> &val) {
  int i, j, v;
  char ch;
  for (i=0; i<freq.size(); i++) {
    for (j=1; j<(freq.size()-i); j++) {
      if (freq[j] > freq[j-1]) {
        v  = freq[j]; freq[j] = freq[j-1]; freq[j-1] = v;
        ch = val[j];  val[j]  = val[j-1];  val[j-1]  = ch;
      }
    }
  }
}

int get_tilestep_from_normalized_name(std::string &name) {
  int i;
  std::string buf;
  long int li;

  for (i=0; i<4; i++) { buf += name[8+i]; }
  li = strtol(buf.c_str(), NULL, 16);

  return (int)li;
}

// Tilespan is in hex
//
int get_tilespan_from_normalized_name(std::string &name) {
  int i;
  std::string buf;
  long int li;

  for (i=0; i<4; i++) { buf += name[18+i]; }
  li = strtol(buf.c_str(), NULL, 16);

  return (int)li;
}

// Create an ASCII canonical sequence from
// the sequecnes from `idx_start` to `idx_n`.
// This will fill-in the sequence based on frequency.
//
void create_canon_seq(std::string &canon_seq, std::vector<sglf2bit_t *> &sglf2bit, size_t idx_start, size_t idx_n, std::vector<std::string> &tagset) {
  int i, j;
  size_t idx, seq_n;
  std::vector<int> freq[4];
  std::vector<char> acgt;
  std::vector<int> pos_freq;
  std::string seq;
  int tilestep_start, tilestep_span;

  canon_seq.clear();

  seq_n = sglf2bit[idx_start]->twobit->size;
  for (j=0; j<4; j++) {
    for (i=0; i<seq_n; i++) {
      freq[j].push_back(0);
    }
  }

  for (idx=idx_start; idx<(idx_start+idx_n); idx++) {
    sglf2bit[idx]->twobit->twoBitToDnaSeq(seq);

    for (i=0; i<seq_n; i++) {

      if      ((seq[i] == 'a') || (seq[i] == 'A')) { freq[0][i]++; }
      else if ((seq[i] == 'c') || (seq[i] == 'C')) { freq[1][i]++; }
      else if ((seq[i] == 'g') || (seq[i] == 'G')) { freq[2][i]++; }
      else if ((seq[i] == 't') || (seq[i] == 'T')) { freq[3][i]++; }
      else { freq[0][i]++; }

    }

  }

  acgt.push_back('a');   acgt.push_back('c');   acgt.push_back('g');   acgt.push_back('t');
  pos_freq.push_back(0); pos_freq.push_back(0); pos_freq.push_back(0); pos_freq.push_back(0);

  for (i=0; i<seq_n; i++) {
    acgt[0] = 'a'; acgt[1] = 'c'; acgt[2] = 'g'; acgt[3] = 't';
    pos_freq[0] = freq[0][i];
    pos_freq[1] = freq[1][i];
    pos_freq[2] = freq[2][i];
    pos_freq[3] = freq[3][i];
    bub_mirror(pos_freq, acgt);
    canon_seq.push_back(acgt[0]);
  }

  tilestep_start = get_tilestep_from_normalized_name(sglf2bit[idx_start]->twobit->name);
  tilestep_span = get_tilespan_from_normalized_name(sglf2bit[idx_start]->twobit->name);

  if (canon_seq.size()<48) { return; }

  // We force tags to be the known tag sequence instead of doing the averaging
  //
  if (tilestep_start > 0) {
    for (i=0; i<24; i++) {
      canon_seq[i] = tagset[tilestep_start-1][i];
    }
  }

  if ((tilestep_start + tilestep_span) <= tagset.size()) {
    for (i=0; i<24; i++) {
      canon_seq[seq_n-24+i] = tagset[tilestep_start + tilestep_span - 1][i];
    }
  }

}

// We use the MD5 hash digest of the sequence
// for uniqueness tests.
//
typedef struct m5_freq_idx_type {
  std::string m5str;
  int freq;
  int idx;
} m5_freq_idx_t;

bool m5_freq_idx_cmp_freq(m5_freq_idx_t &a, m5_freq_idx_t &b) {
  if (b.freq < a.freq) { return true; }
  return false;
}

bool m5_freq_idx_cmp_m5str(m5_freq_idx_t &a, m5_freq_idx_t &b) {
  if (a.m5str < b.m5str) { return true; }
  return false;
}

// Fill in the two bit 'raw' sequence regardless of whether
// it falls within a nocall or other mask.
//
int fillin_raw_seq(twoBit_t *twobit, std::string &canon_seq) {
  size_t i, j, k;
  size_t n, q, r;
  std::string seq;
  unsigned char uc, mask, dna_code;

  n = twobit->size;
  if (n!=canon_seq.size()) { return -1; }

  twobit->twoBitToDnaSeq(seq);

  for (i=0; i<n; i++) {
    if (seq[i] <= 'a') { seq[i] -= 'A'; seq[i] += 'a'; }
    if (seq[i] == 'n') {

      q = i/4;
      r = i%4;

      dna_code = 0;
      uc = twobit->data[q];
      mask = 0xff;
      switch (canon_seq[i]) {
        case 'g': dna_code = 3; break;
        case 'a': dna_code = 2; break;
        case 'c': dna_code = 1; break;
        case 't':
        default:  dna_code = 0; break;
      }

      mask = (0x3 << (2*(3-r)));
      mask = ~mask;

      uc = twobit->data[q] & mask;
      uc |= (dna_code << (2*(3-r)));

      twobit->data[q] = uc;
    }
  }

  return 0;
}

int seq_sanity(std::string &seq) {
  int i;
  for (i=0; i<seq.size(); i++) {
    if ( (seq[i]!='a') && (seq[i]!='c') &&
         (seq[i]!='g') && (seq[i]!='t') ) {
      return -1;
    }
  }
  return 0;
}

// Helper function to create an ASCII representation
// of the MD5 digest from the sequence `seq`
//
void md5str(std::string &s, std::string &seq) {
  int i;
  unsigned char m[MD5_DIGEST_LENGTH];
  char buf[32];

  s.clear();

  MD5((unsigned char *)(seq.c_str()), seq.size(), m);

  for (i=0; i<MD5_DIGEST_LENGTH; i++) {
    sprintf(buf, "%02x", (unsigned char)m[i]);
    s += buf;
  }

}

void print_raw_seq(FILE *fp, twoBit_t *twobit) {
  size_t i, q, r;
  unsigned char uc, dat;
  char c;

  for (i=0; i<twobit->size; i++) {
    c = 'x';
    q = i/4;
    r = i%4;
    dat = twobit->data[q];
    dat >>= 2*(3-r);
    dat &= 0x3;
    if      (dat==0) { c = 't'; }
    else if (dat==1) { c = 'c'; }
    else if (dat==2) { c = 'a'; }
    else if (dat==3) { c = 'g'; }
    fprintf(fp, "%c", c);
  }
}

// Print the slgf2bit array sorted by tile ID and frequency of occurance.
// All sequences printed will have their 'no-call' bases filled in with
// a canonical sequence, which is calculated here.
// The canonical sequence is created by filling in with the most 'frequent'
// base that appears at that sequence position for all sequences that have
// equal size and have the same tile ID.
// The only exception is if a nocall lands on a tag, in which case the tag sequence
// is used.
//
void print_sglf_seq(std::vector<std::string> &tagset, std::vector<sglf2bit_t *> &sglf2bit) {
  int i, j, k, r;
  int idx=0, idx_n=0;
  int seq_idx_s=0, seq_idx_n=0;
  int name_cmp_len = 0;
  int sglf_idx, span;
  std::string canon_seq, raw_seq, m5s;

  std::vector< m5_freq_idx_t > m5_freq_idx_all, m5_freq_idx_dedup;
  m5_freq_idx_t m5fi;

  name_cmp_len = strlen("0000.00.0000.");

  idx=0;
  while (idx < sglf2bit.size()) {

    // first find the block with tiles that have the same tile step
    //
    idx_n=1;
    while ((idx + idx_n) < sglf2bit.size()) {
      if (strncmp(sglf2bit[idx]->twobit->name.c_str(), sglf2bit[idx+idx_n]->twobit->name.c_str(), name_cmp_len)!=0) { break; }
      idx_n++;
    }

    // now find a sub-block within this tile step block where the
    // tile spans and tile sequence lengths are the same.
    // Once we find a group of tiles that have the same tile step, span
    // and sequence size, create the canonical sequence.
    //
    seq_idx_s = idx;
    while (seq_idx_s < (idx + idx_n)) {
      seq_idx_n=1;

      while ((seq_idx_s + seq_idx_n) < (idx + idx_n)) {
        if (sglf2bit[seq_idx_s]->twobit->name != sglf2bit[seq_idx_s + seq_idx_n]->twobit->name) { break; }
        if (sglf2bit[seq_idx_s]->twobit->size != sglf2bit[seq_idx_s + seq_idx_n]->twobit->size) { break; }
        seq_idx_n++;
      }

      create_canon_seq(canon_seq, sglf2bit, seq_idx_s, seq_idx_n, tagset);

      // fill in sequences
      //
      for (k=0; k<seq_idx_n; k++) {
        fillin_raw_seq(sglf2bit[seq_idx_s+k]->twobit, canon_seq);
      }

      seq_idx_s += seq_idx_n;
    }

    // sequences have been filled in, now find frequency of each

    // create array with m5str and indicies
    //
    m5_freq_idx_all.clear();
    for (k=0; k<idx_n; k++) {
      sglf2bit[idx+k]->twobit->twoBitToRawDnaSeq(raw_seq);

      r = seq_sanity(raw_seq);
      if (r<0) {
        fprintf(stderr, "SANITY: sequence has extraneous characters (idx %i)\n",
            (int)(idx+k));
        fflush(stderr);
      }
      md5str(m5s, raw_seq);

      m5fi.freq = sglf2bit[idx+k]->freq;
      m5fi.idx = idx+k;
      m5fi.m5str = m5s;

      m5_freq_idx_all.push_back(m5fi);
    }

    std::sort(m5_freq_idx_all.begin(), m5_freq_idx_all.end(), m5_freq_idx_cmp_m5str);

    // calculate the frequency of each tile
    //
    m5fi.freq  = m5_freq_idx_all[0].freq;
    m5fi.m5str = m5_freq_idx_all[0].m5str;
    m5fi.idx   = m5_freq_idx_all[0].idx;
    m5_freq_idx_dedup.clear();
    for (i=1; i<m5_freq_idx_all.size(); i++) {
      if (m5_freq_idx_all[i-1].m5str == m5_freq_idx_all[i].m5str) {
        m5fi.freq += m5_freq_idx_all[i].freq;
      }
      else {
        m5_freq_idx_dedup.push_back(m5fi);
        m5fi.freq  = m5_freq_idx_all[i].freq;
        m5fi.m5str = m5_freq_idx_all[i].m5str;
        m5fi.idx   = m5_freq_idx_all[i].idx;
      }

    }
    m5_freq_idx_dedup.push_back(m5fi);

    std::sort(m5_freq_idx_dedup.begin(), m5_freq_idx_dedup.end(), m5_freq_idx_cmp_freq);

    // print out SGLF in freq order
    //
    for (i=0; i<m5_freq_idx_dedup.size(); i++) {
      sglf_idx = m5_freq_idx_dedup[i].idx;

      for (j=0; j<12; j++) { printf("%c", sglf2bit[sglf_idx]->twobit->name[j]); }
      span = strtol(sglf2bit[sglf_idx]->twobit->name.c_str() + 18, NULL, 16);
      printf(".%03x+%d,%s,", i, span, m5_freq_idx_dedup[i].m5str.c_str());
      print_raw_seq(stdout, sglf2bit[sglf_idx]->twobit);
      printf("\n");
    }

    idx += idx_n;
  }

}

void cleanup(std::vector<sglf2bit_t *> &sglf2bit) {
  for (int i=0; i<sglf2bit.size(); i++) {
    delete sglf2bit[i]->twobit;
    delete sglf2bit[i];
  }
}

int read_tagset(FILE *fp, std::vector<std::string> &tagset) {
  int ch;
  std::string tag;
  size_t ch_count=0;

  tagset.clear();

  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch=='\n') || (ch==EOF)) { continue; }

    tag += (char)ch;
    if (tag.size()==24) {
      tagset.push_back(tag);
      tag.clear();
    }

  }

  if (tag.size()!=0) { return -1; }
  return 0;
}

int main(int argc, char **argv) {
  int i, j, k;
  std::vector<sglf2bit_t *> sglf2bit;
  std::vector<std::string> tagset;
  FILE *tagset_fp;

  opt_t opt;

  opt.dedup_fj = 1;

  tagset_fp = fopen(argv[1], "r");
  read_tagset(tagset_fp, tagset);
  fclose(tagset_fp);

  read_csvseq_to_twobit(stdin, sglf2bit, opt);
  print_sglf_seq(tagset, sglf2bit);
  cleanup(sglf2bit);
}
