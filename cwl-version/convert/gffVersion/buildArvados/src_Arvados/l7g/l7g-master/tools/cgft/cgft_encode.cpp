#include "cgft.hpp"


// Read in tilemap from tilemap text file
// Text file format is one tilemap. E.g.;
//
// 0:0
// 0:1
// 1:0
// ...
// 4+2:0;0
// ...
//
int load_tilemap(std::string &tilemap_str, std::map< std::string, int > &tilemap) {
  int i, j, k, n;
  int pos=0;
  char buf[1024];
  FILE *fp;
  int idx=0;

  std::map< std::string, int >::iterator ent;
  std::string key;
  std::string item;

  std::stringstream ss;
  ss.str(tilemap_str);

  while (std::getline(ss, item, '\n')) {
    tilemap[item] = pos;
    pos++;
  }

}

void ez_create_enc_vector(sdsl::enc_vector<> &enc_vec, std::vector<int> &v) {
  int i;
  sdsl::int_vector<> iv(v.size());
  for (i=0; i<v.size(); i++) { iv[i] = v[i]; }
  sdsl::enc_vector<> ev(iv);
  enc_vec = ev;
}

void ez_create_vlc_vector(sdsl::vlc_vector<> &vlc_vec, std::vector<int> &v) {
  int i;
  sdsl::int_vector<> iv(v.size());
  for (i=0; i<v.size(); i++) { iv[i] = v[i]; }
  sdsl::vlc_vector<> ve(iv);
  vlc_vec = ve;
}

void ez_to_tilepath(tilepath_t *tilepath, tilepath_ez_t *ez) {
  int i, j, k;
  int val;

  tilepath->NTileStep = (uint64_t)ez->N;

  if (tilepath->Cache) { delete [] tilepath->Cache; }
  tilepath->Cache = new uint64_t[ez->cache.size()];
  for (i=0; i<ez->cache.size(); i++) {
    tilepath->Cache[i] = ez->cache[i];
  }

  if (tilepath->Loq) { delete [] tilepath->Loq; }
  tilepath->Loq = new unsigned char[ez->loq_bv.size()];
  for (i=0; i<ez->loq_bv.size(); i++) {
    tilepath->Loq[i] = ez->loq_bv[i];
  }

  if (tilepath->Span) { delete [] tilepath->Span; }
  tilepath->Span = new unsigned char[ez->span_bv.size()];
  for (i=0; i<ez->span_bv.size(); i++) {
    tilepath->Span[i] = ez->span_bv[i];
  }

  if (tilepath->Overflow) { delete [] tilepath->Overflow; }
  tilepath->Overflow = new uint16_t[ez->ovf_vec.size()];
  for (i=0; i<ez->ovf_vec.size(); i++) {
    tilepath->Overflow[i] = ez->ovf_vec[i];
  }
  tilepath->NOverflow = (uint64_t)ez->ovf_vec.size();

  if (tilepath->Overflow64) { delete tilepath->Overflow64; }
  tilepath->Overflow64 = NULL;
  if (ez->ovf64_vec.size()>0) {
    tilepath->NOverflow64 = (uint64_t)ez->ovf64_vec.size();
    tilepath->Overflow64 = new uint64_t[ez->ovf64_vec.size()];
    for (i=0; i<ez->ovf64_vec.size(); i++) {
      tilepath->Overflow64[i] = ez->ovf64_vec[i];
    }
  }

  if (tilepath->ExtraData) { delete tilepath->ExtraData; }
  tilepath->ExtraData = NULL;
  if (ez->data_vec.size() > 0) {
    tilepath->ExtraDataSize = (uint64_t)ez->data_vec.size();
    tilepath->ExtraData = new char[ez->data_vec.size()];
    for (i=0; i<ez->data_vec.size(); i++) {
      tilepath->ExtraData[i] = ez->data_vec[i];
    }
  }

  std::vector<int> start_v, len_v, sn, variant_v;

  // Encoded Loq Hom information
  //

  ez_create_enc_vector(tilepath->LoqTileStepHom, ez->loq_info_pos_hom);

  // We encode a nocall variant as a positive large number (vlc can't encode
  // negative numbers).
  //

  //ez_create_vlc_vector(tilepath->LoqTileVariantHom, ez->loq_info_variant_hom);
  for (i=0; i<ez->loq_info_variant_hom.size(); i++) {
    val = ez->loq_info_variant_hom[i];
    //if (val<0) { variant_v.push_back(1<<30); }
    if (val<0) { variant_v.push_back(SPAN_SDSL_ENC_VAL); }
    else { variant_v.push_back(val); }
  }
  ez_create_vlc_vector(tilepath->LoqTileVariantHom, variant_v);

  // Sum vector holds sum of the interleaved vector.  Since we split it out,
  // we divide each entry by two to get the entry sum.
  //
  for (i=0; i<ez->loq_info_sn_hom.size(); i++) {
    sn.push_back((int)ez->loq_info_sn_hom[i]/2);
  }
  ez_create_enc_vector(tilepath->LoqTileNocSumHom, sn);

  for (i=0; i<ez->loq_info_noc_hom.size(); i+=2) {
    start_v.push_back(ez->loq_info_noc_hom[i]);
    len_v.push_back(ez->loq_info_noc_hom[i+1]);
  }

  ez_create_vlc_vector(tilepath->LoqTileNocStartHom, start_v);
  ez_create_vlc_vector(tilepath->LoqTileNocLenHom, len_v);


  // Encoded Loq Het information
  //

  start_v.clear();
  len_v.clear();
  sn.clear();
  variant_v.clear();

  ez_create_enc_vector(tilepath->LoqTileStepHet, ez->loq_info_pos_het);
  // We encode a nocall variant as a positive large number (vlc can't encode
  // negative numbers).
  //

  //ez_create_vlc_vector(tilepath->LoqTileVariantHet, ez->loq_info_variant_het);
  for (i=0; i<ez->loq_info_variant_het.size(); i++) {
    val = ez->loq_info_variant_het[i];
    //if (val<0) { variant_v.push_back(1<<30); }
    if (val<0) { variant_v.push_back(SPAN_SDSL_ENC_VAL); }
    else { variant_v.push_back(val); }
  }
  ez_create_vlc_vector(tilepath->LoqTileVariantHet, variant_v);


  // Sum vector holds sum of the interleaved vector.  Since we split it out,
  // we divide each entry by two to get the entry sum.
  //
  for (i=0; i<ez->loq_info_sn_het.size(); i++) {
    sn.push_back((int)ez->loq_info_sn_het[i]/2);
  }
  ez_create_enc_vector(tilepath->LoqTileNocSumHet, sn);

  for (i=0; i<ez->loq_info_noc_het.size(); i+=2) {
    start_v.push_back(ez->loq_info_noc_het[i]);
    len_v.push_back(ez->loq_info_noc_het[i+1]);
  }

  ez_create_vlc_vector(tilepath->LoqTileNocStartHet, start_v);
  ez_create_vlc_vector(tilepath->LoqTileNocLenHet, len_v);

}

int cgft_read_band_tilepath(cgf_t *cgf, tilepath_t *tilepath, FILE *fp) {
  int i, j, k, ch=1;
  int read_line = 0;
  int step=0;

  std::vector<std::string> names;
  std::string s;

  std::vector<tilepath_vec_t> ds;
  tilepath_vec_t cur_ds;

  int pcount=0;
  int state_mod = 0;
  int cur_allele = 0;
  int loq_flag = 0;
  int cur_tilestep = 0;

  const char *fn_tilemap = "default_tile_map_v0.1.0.txt";
  std::map< std::string, int > tilemap;
  std::map< std::string, int >::iterator ent;

  std::vector<int> loq_vec;

  load_tilemap(cgf->TileMap, tilemap);

  s.clear();
  while (ch!=EOF) {
    ch = fgetc(fp);

    if (ch==EOF) { break; }
    if (ch=='\n') {
      state_mod = (state_mod+1)%4;
      pcount=0;
      cur_tilestep=0;

      loq_vec.clear();

      if (state_mod==0) {
        cur_tilestep=0;
        ds.push_back(cur_ds);
        cur_ds.allele[0].clear();
        cur_ds.allele[1].clear();
        cur_ds.loq_flag[0].clear();
        cur_ds.loq_flag[1].clear();
        cur_ds.loq_info[0].clear();
        cur_ds.loq_info[1].clear();
        cur_ds.name.clear();
      }

      continue;
    }
    if (ch=='[') {

      loq_flag=0;
      if (state_mod>=2) {

        s.clear();
        loq_vec.clear();

        pcount++;
        while (pcount>1) {
          cur_allele = state_mod%2;
          ch = fgetc(fp);


          if (ch==EOF) {
            printf("ERROR: premature eof\n");
            return -1;
            //exit(1);
          }

          if ((ch==' ') || (ch==']')) {
            if (s.size() > 0) {
              loq_vec.push_back(atoi(s.c_str()));
            }
            s.clear();
          }


          if (ch==']') {
            cur_ds.loq_flag[cur_allele].push_back(loq_flag);
            pcount--;

            cur_ds.loq_info[cur_allele].push_back(loq_vec);
            loq_vec.clear();
            cur_tilestep++;
            continue;
          }
          if (ch=='[') { pcount++; continue; }
          if (ch==' ') { continue; }

          s += ch;
          loq_flag=1;
        }
      }

      continue;
    }
    if ((ch==' ') || (ch==']')) {

      if (s.size() == 0) { continue; }

      if (state_mod==0) {
        cur_ds.allele[0].push_back(atoi(s.c_str()));
      } else if (state_mod==1) {
        cur_ds.allele[1].push_back(atoi(s.c_str()));
      }
      s.clear();
      continue;
    }
    s += ch;

  }

  //DEBUG
  //print_bgf(ds[0]);
  //print_tilepath_vec(ds[0]);

  tilepath_ez_t ez;
  ez_create(ez, ds[0], tilemap);

  //DEBUG
  //ez_print(ez);

  /*
  for (i=0; i<ez.loq_info_pos.size(); i++) {
    printf(" ez loq[%i]: variant %i: %i %i\n",
        i,
        ez.loq_info_pos[i],
        ez.loq_info_variant[2*i],
        ez.loq_info_variant[2*i+1]);
  }

  printf("HOM:\n");

  for (i=0; i<ez.loq_info_pos_hom.size(); i++) {
    printf(" ez hom loq[%i]: variant %i: %i %i\n",
        i,
        ez.loq_info_pos_hom[i],
        ez.loq_info_variant_hom[2*i],
        ez.loq_info_variant_hom[2*i+1]);
  }

  printf("HET:\n");

  for (i=0; i<ez.loq_info_pos_het.size(); i++) {
    printf(" ez het loq[%i]: variant %i: %i %i\n",
        i,
        ez.loq_info_pos_het[i],
        ez.loq_info_variant_het[2*i],
        ez.loq_info_variant_het[2*i+1]);
  }
  */

  ez_to_tilepath(tilepath, &ez);


  return 0;

}
