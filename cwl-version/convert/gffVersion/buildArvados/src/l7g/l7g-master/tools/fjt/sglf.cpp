#include "sglf.hpp"

int sglf_path_step_lookup_hash_variant_id(sglf_path_t *sp, int tilestep, std::string &hash) {
  int var_idx=-1;

  if ((tilestep < 0) || (tilestep > ((int)(sp->hash.size())))) { return -1; }

  for (var_idx=0; var_idx<(sp->hash[tilestep].size()); var_idx++) {
    if ( sp->hash[tilestep][var_idx] == hash ) { return var_idx; }
  }

  return -1;
}

int sglf_path_step_lookup_seq_variant_id(sglf_path_t *sp, int tilestep, std::string &seq) {
  int var_idx=-1, i;
  size_t n;

  if ((tilestep < 0) || (tilestep > ((int)(sp->seq.size())))) { return -1; }

  for (var_idx=0; var_idx<(sp->seq[tilestep].size()); var_idx++) {

    //DEBUG
    //printf("sglf_p_s_seq_var_lkp: varidx: %i, sz[%i] %i, seq %i\n",
    //    var_idx,
    //    tilestep,
    //    (int)(sp->seq[tilestep].size()),
    //    (int)seq.size());

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

