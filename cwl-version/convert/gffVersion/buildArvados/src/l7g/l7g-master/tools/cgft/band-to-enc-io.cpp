#include "cgft.hpp"

using namespace sdsl;

int ez_load(const char *base_f, cgf_ez_t &cgf_ez) {
  return 0;
}

int ez_save_vec64(const char *fn, std::vector<uint64_t> &v) {
  FILE *ofp;
  ofp = fopen(fn, "w");
  if (!ofp) { return -1; }
  fwrite(&(v[0]), sizeof(uint64_t), v.size(), ofp);
  fclose(ofp);
  return 0;
}

int ez_save_vec16(const char *fn, std::vector<int16_t> &v) {
  FILE *ofp;
  ofp = fopen(fn, "w");
  if (!ofp) { return -1; }
  fwrite(&(v[0]), sizeof(int16_t), v.size(), ofp);
  fclose(ofp);
  return 0;
}

int ez_save_vec8(const char *fn, std::vector<unsigned char> &v) {
  FILE *ofp;
  ofp = fopen(fn, "w");
  if (!ofp) { return -1; }
  fwrite(&(v[0]), sizeof(unsigned char), v.size(), ofp);
  fclose(ofp);
  return 0;
}

int ez_save_sdsl_ivec(const char *fn, std::vector<int> &v) {
  int i;
  int_vector<> iv(v.size());
  for (i=0; i<v.size(); i++) { iv[i] = v[i]; }
  vlc_vector<> ve(iv);
  store_to_file(ve, fn);
  return 0;
}

int ez_save_sdsl_ivec_enc(const char *fn, std::vector<int> &v) {
  int i;
  int_vector<> iv(v.size());
  for (i=0; i<v.size(); i++) { iv[i] = v[i]; }
  enc_vector<> ev(iv);
  store_to_file(ev, fn);
  return 0;
}

int ez_save_sdsl_ivec_alt(const char *fn, std::vector<int> &v) {
  int i;
  int s=0;

  for (i=0; i<v.size(); i++) { s+=v[i]; }

  bit_vector b = bit_vector(s, 0);

  for (i=0; i<v.size(); i++)  { b[i] = 1; }

  sd_vector<> sdb(b);
  store_to_file(sdb, fn);

  return 0;
}

int ez_save(const char *base_f, int tilepath, tilepath_ez_t &ez) {
  int i;
  FILE *fp;
  char buf[1024], hex_tps[1024];
  std::string fn, base_ofn;

  std::vector<int> x,y;

  sprintf(hex_tps, "%03x", tilepath);

  base_ofn.clear();
  base_ofn += base_f;
  base_ofn += "/";
  base_ofn += hex_tps;

  fn = base_ofn; fn += "-cache";
  ez_save_vec64(fn.c_str(), ez.cache);

  fn = base_ofn; fn += "-span_bv";
  ez_save_vec8(fn.c_str(), ez.span_bv);

  fn = base_ofn; fn += "-ovf_vec";
  ez_save_vec16(fn.c_str(), ez.ovf_vec);

  fn = base_ofn; fn += "-loq_vec";
  ez_save_vec8(fn.c_str(), ez.loq_bv);

  /*
  fn = base_ofn; fn += "-loq_info_pos.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_pos);

  //clean up variants
  for (i=0; i<ez.loq_info_variant.size(); i++) {
    if (ez.loq_info_variant[i]<0) {
      ez.loq_info_variant[i] = 1<<20;
    }
  }

  fn = base_ofn; fn += "-loq_info_variant.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_variant);

  fn = base_ofn; fn += "-loq_info_sn.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_sn);

  //fn = base_ofn;
  //fn += "-loq_info_sn-sdb.sdsl";
  //ez_save_sdsl_ivec_alt(fn.c_str(), ez.loq_info_sn);

  //fn = base_ofn;
  //fn += "-loq_info_noc.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_noc);

  for (i=0; i<ez.loq_info_noc.size(); i+=2) {
    x.push_back(ez.loq_info_noc[i]);
    y.push_back(ez.loq_info_noc[i+1]);
  }

  fn = base_ofn; fn += "-loq_info_noc-start.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), x);

  fn = base_ofn; fn += "-loq_info_noc-len.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), y);
  */

  //--------------------

  // hom loq info
  //
  fn = base_ofn; fn += "-loq_info_pos_hom.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_pos_hom);
  ez_save_sdsl_ivec_enc(fn.c_str(), ez.loq_info_pos_hom);

  //clean up variants
  //
  for (i=0; i<ez.loq_info_variant_hom.size(); i++) {
    if (ez.loq_info_variant_hom[i]<0) {
      ez.loq_info_variant_hom[i] = 1<<20;
    }
  }

  fn = base_ofn; fn += "-loq_info_variant_hom.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_variant_hom);

  fn = base_ofn; fn += "-loq_info_sn_hom.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_sn_hom);
  ez_save_sdsl_ivec_enc(fn.c_str(), ez.loq_info_sn_hom);

  //fn = base_ofn; fn += "-loq_info_noc_hom.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_noc_hom);

  x.clear(); y.clear();
  for (i=0; i<ez.loq_info_noc_hom.size(); i+=2) {
    x.push_back(ez.loq_info_noc_hom[i]);
    y.push_back(ez.loq_info_noc_hom[i+1]);
  }

  fn = base_ofn; fn += "-loq_info_noc_hom-start.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), x);
  //ez_save_sdsl_ivec_enc(fn.c_str(), x);

  fn = base_ofn; fn += "-loq_info_noc_hom-len.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), y);
  //ez_save_sdsl_ivec_enc(fn.c_str(), y);


  //--------------------

  // het loq info
  //
  fn = base_ofn; fn += "-loq_info_pos_het.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_pos_het);
  ez_save_sdsl_ivec_enc(fn.c_str(), ez.loq_info_pos_het);

  //clean up variants
  //
  for (i=0; i<ez.loq_info_variant_het.size(); i++) {
    if (ez.loq_info_variant_het[i]<0) {
      ez.loq_info_variant_het[i] = 1<<20;
    }
  }

  fn = base_ofn; fn += "-loq_info_variant_het.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_variant_het);

  fn = base_ofn; fn += "-loq_info_sn_het.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_sn_het);
  ez_save_sdsl_ivec_enc(fn.c_str(), ez.loq_info_sn_het);

  //fn = base_ofn; fn += "-loq_info_noc_het.sdsl";
  //ez_save_sdsl_ivec(fn.c_str(), ez.loq_info_noc_het);

  x.clear(); y.clear();
  for (i=0; i<ez.loq_info_noc_het.size(); i+=2) {
    x.push_back(ez.loq_info_noc_het[i]);
    y.push_back(ez.loq_info_noc_het[i+1]);
  }

  fn = base_ofn; fn += "-loq_info_noc_het-start.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), x);
  //ez_save_sdsl_ivec_enc(fn.c_str(), x);

  fn = base_ofn; fn += "-loq_info_noc_het-len.sdsl";
  ez_save_sdsl_ivec(fn.c_str(), y);
  //ez_save_sdsl_ivec_enc(fn.c_str(), y);


  return 0;
}
