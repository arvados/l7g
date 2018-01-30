#include "sglf.hpp"

int sglf_path_step_lookup_hash_variant_id(sglf_path_t *sp, int tilestep, std::string &hash) {
  int var_idx=-1;

  if ((tilestep < 0) || (tilestep > ((int)(sp->hash.size())))) { return -1; }

  for (var_idx=0; var_idx<(sp->hash[tilestep].size()); var_idx++) {
    if ( sp->hash[tilestep][var_idx] == hash ) { return var_idx; }
  }

  return -1;
}

int sglf2bit_tilepath_step_lookup_seq_variant_id(sglf2bit_tilepath_t &sp, int tilestep, std::string &seq) {
  int var_idx=-1, i;
  size_t n;
  std::string sglf_seq;

  if ((tilestep < 0) || (tilestep > ((int)(sp.seq2bit.size())))) { return -1; }

  for (var_idx=0; var_idx<(sp.seq2bit[tilestep].size()); var_idx++) {

    if (sp.seq2bit[tilestep][var_idx].size != seq.size()) { continue; }
    n = seq.size();

    sp.seq2bit[tilestep][var_idx].twoBitToDnaSeq(sglf_seq);

    for (i=0; i<n; i++) {

      if ((seq[i] == 'n') || (seq[i] == 'N')) { continue; }
      if (seq[i] != sglf_seq[i]) { break; }

    }
    if (i==n) { return var_idx; }

  }

  return -2;
}

int sglf_path_step_lookup_seq_variant_id(sglf_path_t *sp, int tilestep, std::string &seq) {
  int var_idx=-1, i;
  size_t n;

  if ((tilestep < 0) || (tilestep > ((int)(sp->seq.size())))) { return -1; }

  for (var_idx=0; var_idx<(sp->seq[tilestep].size()); var_idx++) {

    if (sp->seq[tilestep][var_idx].size() != seq.size()) { continue; }
    n = seq.size();

    for (i=0; i<n; i++) {

      if ((seq[i] == 'n') || (seq[i] == 'N')) { continue; }
      if (seq[i] != sp->seq[tilestep][var_idx][i]) { break; }

    }
    if (i==n) { return var_idx; }

  }

  return -2;
}


void sglf_path_print(sglf_path_t *sp) {
  int tilestep, tilevar;
  int n, m;

  n = (int)(sp->ext_tileid.size());

  for (tilestep=0; tilestep<n; tilestep++) {
    m = (int)(sp->ext_tileid[tilestep].size());

    for (tilevar=0; tilevar<m; tilevar++) {
      printf("%s,%s,%s\n",
          sp->ext_tileid[tilestep][tilevar].c_str(),
          sp->hash[tilestep][tilevar].c_str(),
          sp->seq[tilestep][tilevar].c_str());
    }

  }

}


uint16_t tileid_part(uint64_t tileid, int part) {
  uint64_t u64;
  uint16_t u16;
  u64 = tileid>>(8*2*part);
  u64 &= 0xffff;
  u16 = (uint16_t)u64;
  return u16;
}


uint64_t parse_tileid(const char *tileid) {
  const char *chp;
  std::string s;
  unsigned long long int ull;
  uint64_t v=0, u64;
  int curpos=0;
  unsigned int byte_offset[] = { 6, 4, 2, 0 };

  for (chp=tileid; *chp; chp++) {
    //if (*chp == '.') {
    if ((*chp == '.') || (*chp=='+')) {
      ull = strtoull(s.c_str(), NULL, 16);
      u64 = (uint64_t)ull;
      v |= (u64 << (8*byte_offset[curpos]));
      curpos++;
      if (curpos>=4) { break; }
      s.clear();
      continue;
    }

    s+=*chp;
  }

  if (curpos<4) {
    ull = strtoull(s.c_str(), NULL, 16);
    u64 = (uint64_t)ull;
    v |= (u64 << (8*byte_offset[curpos]));
  }

  return v;
}


int read_sglf_path(FILE *ifp, sglf_path_t &sp) {
  int i;
  int ch;
  std::string line;
  std::string tid_str, hash_str, seq;
  int cur_tilestep=0, cur_tilevar_idx=0;
  int prev_tilestep=-1, prev_tilevar_idx=-1;
  uint64_t tileid;
  std::pair< int, int > ipair;

  int read_state = 0;

  std::vector< std::string > svec;

  sp.ext_tileid.clear();
  sp.hash.clear();
  sp.seq.clear();

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {

      tileid = parse_tileid(tid_str.c_str());
      cur_tilestep = (int)tileid_part(tileid, 1);
      cur_tilevar_idx = (int)tileid_part(tileid, 0);

      if ( cur_tilestep >= (int)(sp.ext_tileid.size()) ) {
        int del = cur_tilestep - (int)(sp.ext_tileid.size()) + 1;
        for (i=0; i<del; i++) {

          svec.clear();
          sp.ext_tileid.push_back(svec);
          sp.hash.push_back(svec);
          sp.seq.push_back(svec);
        }
      }

      sp.ext_tileid[cur_tilestep].push_back(tid_str);
      sp.hash[cur_tilestep].push_back(hash_str);
      sp.seq[cur_tilestep].push_back(seq);

      ipair.first = cur_tilestep;
      ipair.second = cur_tilevar_idx;
      sp.hash_pos[hash_str] = ipair;

      read_state=0;
      tid_str.clear();
      hash_str.clear();
      seq.clear();

      continue;
    }

    if (ch==',') {
      read_state++;
      continue;
    }

    if (read_state==0)      { tid_str += ch; }
    else if (read_state==1) { hash_str += ch; }
    else if (read_state==2) { seq += ch; }

  }

  //TODO reorder based on pos_map


  return 0;
}


int read_sglf2bit_tilepath(FILE *ifp, sglf2bit_tilepath_t &sp) {
  int i;
  int ch;
  std::string line;
  std::string tid_str, hash_str, seq;
  int cur_tilestep=0, cur_tilevar_idx=0;
  int prev_tilestep=-1, prev_tilevar_idx=-1;
  uint64_t tileid;
  std::pair< int, int > ipair;
  twoBit_t twobit;
  std::vector< twoBit_t > twobit_v;

  int read_state = 0;

  std::vector< std::string > svec;

  sp.ext_tileid.clear();
  sp.hash.clear();
  sp.seq2bit.clear();
  twobit.clear();

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {

      tileid = parse_tileid(tid_str.c_str());
      cur_tilestep = (int)tileid_part(tileid, 1);
      cur_tilevar_idx = (int)tileid_part(tileid, 0);

      if ( cur_tilestep >= (int)(sp.ext_tileid.size()) ) {
        int del = cur_tilestep - (int)(sp.ext_tileid.size()) + 1;
        for (i=0; i<del; i++) {

          twobit.clear();
          svec.clear();
          sp.ext_tileid.push_back(svec);
          sp.hash.push_back(svec);
          sp.seq2bit.push_back(twobit_v);
        }
      }

      sp.ext_tileid[cur_tilestep].push_back(tid_str);
      sp.hash[cur_tilestep].push_back(hash_str);

      twobit.twoBitFromDnaSeq( seq.c_str() );
      sp.seq2bit[cur_tilestep].push_back(twobit);

      ipair.first = cur_tilestep;
      ipair.second = cur_tilevar_idx;
      sp.hash_pos[hash_str] = ipair;

      read_state=0;
      tid_str.clear();
      hash_str.clear();
      seq.clear();
      twobit.clear();

      continue;
    }

    if (ch==',') {
      read_state++;
      continue;
    }

    if (read_state==0)      { tid_str += ch; }
    else if (read_state==1) { hash_str += ch; }
    else if (read_state==2) { seq += ch; }

  }

  return 0;
}
