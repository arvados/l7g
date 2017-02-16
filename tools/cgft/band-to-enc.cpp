// convert output of band cgb file (ascii text) to
// cgf (v4?).
// Still experimental.
//

#include "cgft.hpp"

// We need to experiment but the basic format is something like:
//
// num tot:       N
// span vector:   [ bit vector ]
// cache vector:  [ 32 canon bits | 8 * 4 cache hexits ]
// ovf vector:    [ pos | a | b ]
// data vector:   [ ... ]
//
// num loq:       m
// loq vector:    [ bit vector ]
// loq info:      [ pos | sn_a | sn_b | [ start | run  ... ] | [ start | run ... ] ]
//
// The bit vectors are straight bit vectors, the cache vector is the custom format,
// the ovf vector is a simple array, the data vector is arbitrary data, the loq info
// is an sdsl-lite data structure of packed data, maybe split out by each field to
// facilitate compression.  The 'sn_(a|b)' are the running sums of the number of
// (start, run) pairs so that we keep random access into the (start, run) vector.
//
// There should be some other meta information as well, including the tile map.
// We should try and use as much as the old infrastructure/format as possible.
//

int TILEPATH;

void ez_print(tilepath_ez_t &ez) {
  int s=0;
  int i, j;
  unsigned char u8, u4;
  uint32_t u32;
  uint64_t u64;
  int count;


  printf("ez.N: %i\n", ez.N);

  printf("ez.cache(%i):", (int)(ez.cache.size()));
  for (i=0; i<ez.cache.size(); i++) {
    printf("\n");
    //if ((i>0) && ((i%8)==0)) { printf("\n"); }

    u64 = ez.cache[i];

    u32 = (uint32_t)(u64 >> 32);
    printf(" %08x", u32);

    count=0;
    for (j=0; j<32; j++) {
      if (u32 & (1u<<j)) { count++; }
    }
    if (count>8) { count=8; }

    //force
    count=8;

    printf(" [");
    u32 = (uint32_t)(u64 & 0xfffffffful);
    for (j=0; j<count; j++) {
      u4 = (u32 >> (4*j)) & (0xf);
      if (u4>0) {
        printf(" %01x", u4);
      } else {
        printf(" .");
      }
    }
    printf(" ]");


  }
  printf("\n");

  printf("ez.span_bv(%i):", (int)(ez.span_bv.size()));
  for (i=0; i<ez.span_bv.size(); i++) {
    printf(" %02x", ez.span_bv[i]);
  }
  printf("\n");

  printf("ez.ovf_vec(%i):", (int)(ez.ovf_vec.size()));
  for (i=0; i<ez.ovf_vec.size(); i+=3) {
    printf(" (%i %i %i)", ez.ovf_vec[i], ez.ovf_vec[i+1], ez.ovf_vec[i+2]);
  }
  printf("\n");

  printf("\n");
  printf("ez.loq_bv(%i):", (int)(ez.loq_bv.size()));
  for (i=0; i<ez.loq_bv.size(); i++) {
    printf(" %02x", ez.loq_bv[i]);
  }
  printf("\n");

  printf("ez.loq lens %i %i %i %i\n",
      (int)ez.loq_info_pos.size(),
      (int)ez.loq_info_variant.size(),
      (int)ez.loq_info_sn.size(),
      (int)ez.loq_info_noc.size());
  printf("ez.loq(%i):", (int)(ez.loq_info_pos.size()));
  for (i=0; i<ez.loq_info_pos.size(); i++) {
    printf(" {%i(%i,%i)[%i,%i]:",
        ez.loq_info_pos[i],
        ez.loq_info_variant[2*i], ez.loq_info_variant[2*i+1],
        ez.loq_info_sn[2*i], ez.loq_info_sn[2*i+1]);

    printf("[");
    for (j=s; j<ez.loq_info_sn[2*i]; j+=2) {
      if (j>s) { printf(", "); }
      printf("%i+%i", ez.loq_info_noc[j], ez.loq_info_noc[j+1]);
    }
    printf("]");

    s = ez.loq_info_sn[2*i];

    printf(" [");
    for (j=s; j<ez.loq_info_sn[2*i+1]; j+=2) {
      if (j>s) { printf(", "); }
      printf("%i+%i", ez.loq_info_noc[j], ez.loq_info_noc[j+1]);
    }
    printf("]");

    s = ez.loq_info_sn[2*i+1];

    printf("}");
  }
  printf("\n");


  //DEBUG
  /*

  printf("\n");

  printf("pos_hom(%i):", (int)ez.loq_info_pos_hom.size());
  for (i=0; i<ez.loq_info_pos_hom.size(); i++) { printf(" %i", ez.loq_info_pos_hom[i]); }
  printf("\n");

  printf("variant_hom(%i):", (int)ez.loq_info_variant_hom.size());
  for (i=0; i<ez.loq_info_variant_hom.size(); i++) { printf(" %i", ez.loq_info_variant_hom[i]); }
  printf("\n");

  printf("sn_hom(%i):", (int)ez.loq_info_sn_hom.size());
  for (i=0; i<ez.loq_info_sn_hom.size(); i++) { printf(" %i", ez.loq_info_sn_hom[i]); }
  printf("\n");

  printf("noc_hom(%i):", (int)ez.loq_info_noc_hom.size());
  for (i=0; i<ez.loq_info_noc_hom.size(); i++) { printf(" %i", ez.loq_info_noc_hom[i]); }
  printf("\n");
  printf("\n");

  */
  //DEBUG

  printf("ez.loq_hom(%i):", (int)(ez.loq_info_pos_hom.size()));

  s=0;
  for (i=0; i<ez.loq_info_pos_hom.size(); i++) {

    printf(" {%i(%i,%i)[%i]:",
        ez.loq_info_pos_hom[i],
        ez.loq_info_variant_hom[2*i], ez.loq_info_variant_hom[2*i+1],
        ez.loq_info_sn_hom[i]);

    printf("[");
    for (j=s; j<ez.loq_info_sn_hom[i]; j+=2) {
      if (j>s) { printf(", "); }
      printf("%i+%i", ez.loq_info_noc_hom[j], ez.loq_info_noc_hom[j+1]);
    }
    printf("]");

    s = ez.loq_info_sn_hom[i];

    printf("}");
  }
  printf("\n");

  printf("ez.loq_het(%i):", (int)(ez.loq_info_pos_het.size()));

  s=0;
  for (i=0; i<ez.loq_info_pos_het.size(); i++) {

    printf(" {%i(%i,%i)[%i,%i]:",
        ez.loq_info_pos_het[i],
        ez.loq_info_variant_het[2*i], ez.loq_info_variant_het[2*i+1],
        ez.loq_info_sn_het[2*i], ez.loq_info_sn_het[2*i+1]);

    printf("[");
    for (j=s; j<ez.loq_info_sn_het[2*i]; j+=2) {
      if (j>s) { printf(", "); }
      printf("%i+%i", ez.loq_info_noc_het[j], ez.loq_info_noc_het[j+1]);
    }
    printf("]");

    s = ez.loq_info_sn_het[2*i];

    printf("[");
    for (j=s; j<ez.loq_info_sn_het[2*i+1]; j+=2) {
      if (j>s) { printf(", "); }
      printf("%i+%i", ez.loq_info_noc_het[j], ez.loq_info_noc_het[j+1]);
    }
    printf("]");

    printf("}");
  }
  printf("\n");

}

void mk_tilemap_key(std::string &key, tilepath_vec_t &tilepath, int tilestep, int n) {
  int i, j, k, a;
  int prev_step, cur_step;
  char buf[1024];

  key.clear();

  for (a=0; a<2; a++) {

    if (a>0) { key += ":"; }

    prev_step = tilestep;
    for (cur_step=tilestep; cur_step < (tilestep+n); cur_step++) {
      if (tilepath.allele[a][cur_step] < 0) { continue; }
      if (cur_step > tilestep) { key += ";"; }

      sprintf(buf, "%x", tilepath.allele[a][cur_step]);
      key += buf;

      if ((cur_step - prev_step) > 1) {
        sprintf(buf, "+%x", cur_step-prev_step);
        key += buf;
      }
      prev_step = cur_step;
    }

    if ((cur_step - prev_step)>1) {
      sprintf(buf, "+%x", cur_step-prev_step);
      key += buf;
    }

  }

  return;
}

void ez_create(tilepath_ez_t &ez, tilepath_vec_t &tilepath, std::map< std::string, int > &tilemap) {
  int i, j, k, n, m;
  int CACHE_N=32;
  //int CACHE_N=24;
  int HEXIT_N;
  uint64_t u64;
  uint32_t u32, u32_canon, u32_hexit;
  unsigned char u8;

  std::vector<int> hexit_vec;
  std::vector< std::vector< std::vector<int> > > knot_vec;

  int ii, jj;
  int tilestep=0;
  int n_q, n_r;

  int loq_sn = 0;
  int loq_sn_hom = 0, loq_sn_het = 0;

  //int st, en;
  int block_start, en;
  std::vector<int> loq_flag, span_flag, anchor_flag;
  std::string str;

  char buf[1024];

  std::map< std::string, int >::iterator tilemap_it;

  int loc_debug = 0;

  for (i=0; i<1024; i++) { buf[i] = '\0'; }

  n = tilepath.allele[0].size();

  HEXIT_N = (64-CACHE_N)/4;

  ez.N = n;
  ez.cache.clear();
  ez.span_bv.clear();
  ez.ovf_vec.clear();
  ez.ovf32_vec.clear();
  ez.data_vec.clear();
  ez.loq_bv.clear();
  ez.loq_info_pos.clear();
  ez.loq_info_variant.clear();
  ez.loq_info_sn.clear();
  ez.loq_info_noc.clear();

  ez.loq_info_pos_hom.clear();
  ez.loq_info_variant_hom.clear();
  ez.loq_info_sn_hom.clear();
  ez.loq_info_noc_hom.clear();

  ez.loq_info_pos_het.clear();
  ez.loq_info_variant_het.clear();
  ez.loq_info_sn_het.clear();
  ez.loq_info_noc_het.clear();


  // fill out span and loq knot flags
  //
  for (i=0; i<n; i++) {
    loq_flag.push_back(0);
    span_flag.push_back(0);
    anchor_flag.push_back(0);
  }

  // fill out span info first for use later
  //
  for (tilestep=0; tilestep<n; tilestep++) {

    if ((tilepath.allele[0][tilestep]<0) || (tilepath.allele[1][tilestep]<0)) {

      for (block_start=tilestep; block_start>0; block_start--) {
        if ((tilepath.allele[0][block_start]>=0) && (tilepath.allele[1][block_start]>=0)) { break; }
      }

      for (en=tilestep; en<n; en++) {
        if ((tilepath.allele[0][en]>=0) && (tilepath.allele[1][en]>=0)) { break; }
      }

      for (i=block_start; i<en; i++) {
        span_flag[i] = block_start;
      }

      anchor_flag[block_start] = en-block_start;

      tilestep = en-1;
      continue;

    }

    span_flag[tilestep] = tilestep;
  }

  // use span information to fill out nocall runs
  //
  for (tilestep=0; tilestep<n; tilestep++) {
    if (tilepath.loq_flag[0][tilestep] || tilepath.loq_flag[1][tilestep]) {
      block_start = span_flag[tilestep];
      for (i=tilestep; i>=block_start; i--) { loq_flag[i] = 1; }
      for (i=tilestep; i<n; i++) {
        if (span_flag[i] != block_start) { break; }
        loq_flag[i] = 1;
      }
    }
  }

  //DEBUG
  //for (i=0; i<n; i++) { printf("[%i] %i anchor:%i (loq:%i)\n", i, span_flag[i], anchor_flag[i], loq_flag[i]); }


  // -----------------
  //
  // loq_bv
  // loq_info_(pos|sn|noc)
  //

  // Fist populate loq_vec
  //
  n_q = n / 8;
  n_r = n % 8;

  tilestep=0;
  for (ii=0; ii<n_q; ii++) {
    u8=0;
    for (i=0; i<8; i++) {

      // We only store actual tile nocall information rather than worrying
      // about the nocall run
      //
      //if (tilepath.loq_flag[0][tilestep] || tilepath.loq_flag[1][tilestep]) {
      if (loq_flag[tilestep]) {

        int is_het = 0;
        if (tilepath.loq_info[0][tilestep].size() == tilepath.loq_info[1][tilestep].size()) {
          for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
            if ((tilepath.loq_info[0][tilestep][j] != tilepath.loq_info[1][tilestep][j]) ||
                (tilepath.loq_info[0][tilestep][j+1] != tilepath.loq_info[1][tilestep][j+1])) {
              is_het = 1;
              break;
            }
          }
        } else {
          is_het = 1;
        }

        //DEBUG
        /*
        printf(">> tilestep[%i] loq(%i,%i) is_het %i, var(%i,%i)\n",
            tilestep,
            tilepath.loq_flag[0][tilestep],
            tilepath.loq_flag[1][tilestep],
            is_het,
            tilepath.allele[0][tilestep],
            tilepath.allele[1][tilestep]);
            */

        if (!is_het) {

          // loq hom info
          //

          ez.loq_info_pos_hom.push_back(tilestep);

          // note, the loq is het or hom, not the tile variants
          //
          ez.loq_info_variant_hom.push_back(tilepath.allele[0][tilestep]);
          ez.loq_info_variant_hom.push_back(tilepath.allele[1][tilestep]);

          loq_sn_hom += tilepath.loq_info[0][tilestep].size();
          ez.loq_info_sn_hom.push_back(loq_sn_hom);
          for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
            ez.loq_info_noc_hom.push_back(tilepath.loq_info[0][tilestep][j]);
            ez.loq_info_noc_hom.push_back(tilepath.loq_info[0][tilestep][j+1]);
          }

        } else {

          // loq het info
          //
          ez.loq_info_pos_het.push_back(tilestep);
          ez.loq_info_variant_het.push_back(tilepath.allele[0][tilestep]);
          loq_sn_het += (int)(tilepath.loq_info[0][tilestep].size());
          ez.loq_info_sn_het.push_back(loq_sn_het);
          for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
            ez.loq_info_noc_het.push_back(tilepath.loq_info[0][tilestep][j]);
            ez.loq_info_noc_het.push_back(tilepath.loq_info[0][tilestep][j+1]);
          }

          ez.loq_info_variant_het.push_back(tilepath.allele[1][tilestep]);
          loq_sn_het += (int)(tilepath.loq_info[1][tilestep].size());
          ez.loq_info_sn_het.push_back(loq_sn_het);
          for (j=0; j<tilepath.loq_info[1][tilestep].size(); j+=2) {
            ez.loq_info_noc_het.push_back(tilepath.loq_info[1][tilestep][j]);
            ez.loq_info_noc_het.push_back(tilepath.loq_info[1][tilestep][j+1]);
          }

        }


        // "old" duplicated information
        //
        ez.loq_info_pos.push_back(tilestep);
        ez.loq_info_variant.push_back(tilepath.allele[0][tilestep]);
        loq_sn += (int)(tilepath.loq_info[0][tilestep].size());
        ez.loq_info_sn.push_back(loq_sn);
        for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
          ez.loq_info_noc.push_back(tilepath.loq_info[0][tilestep][j]);
          ez.loq_info_noc.push_back(tilepath.loq_info[0][tilestep][j+1]);
        }

        ez.loq_info_variant.push_back(tilepath.allele[1][tilestep]);
        loq_sn += (int)(tilepath.loq_info[1][tilestep].size());
        ez.loq_info_sn.push_back(loq_sn);
        for (j=0; j<tilepath.loq_info[1][tilestep].size(); j+=2) {
          ez.loq_info_noc.push_back(tilepath.loq_info[1][tilestep][j]);
          ez.loq_info_noc.push_back(tilepath.loq_info[1][tilestep][j+1]);
        }

      }

      // For the loq_bv, we store the actual low quality information
      // on the 'knot' level
      //
      if (loq_flag[tilestep]) { u8 |= 1<<i; }

      tilestep++;

    }
    ez.loq_bv.push_back(u8);
  }

  u8=0;
  for (i=0; i<n_r; i++) {

    //if (tilepath.loq_flag[0][tilestep] || tilepath.loq_flag[1][tilestep]) {
    if (loq_flag[tilestep]) {
      u8 |= 1<<i;

      //if (tilepath.loq_info[0][tilestep].size() == tilepath.loq_info[1][tilestep].size()) {

      int is_het = 0;
      if (tilepath.loq_info[0][tilestep].size() == tilepath.loq_info[1][tilestep].size()) {
        for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
          if ((tilepath.loq_info[0][tilestep][j] != tilepath.loq_info[1][tilestep][j]) ||
              (tilepath.loq_info[0][tilestep][j+1] != tilepath.loq_info[1][tilestep][j+1])) {
            is_het = 1;
            break;
          }
        }
      } else {
        is_het = 1;
      }

      if (!is_het) {


        // loq hom info
        //

        ez.loq_info_pos_hom.push_back(tilestep);

        // note, the loq is het or hom, not the tile variants
        //
        ez.loq_info_variant_hom.push_back(tilepath.allele[0][tilestep]);
        ez.loq_info_variant_hom.push_back(tilepath.allele[1][tilestep]);

        loq_sn_hom += tilepath.loq_info[0][tilestep].size();
        ez.loq_info_sn_hom.push_back(loq_sn_hom);
        for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
          ez.loq_info_noc_hom.push_back(tilepath.loq_info[0][tilestep][j]);
          ez.loq_info_noc_hom.push_back(tilepath.loq_info[0][tilestep][j+1]);
        }

      } else {

        //DEBUG
        printf(">>>>>>>\n");

        // loq het info
        //
        ez.loq_info_pos_het.push_back(tilestep);
        ez.loq_info_variant_het.push_back(tilepath.allele[0][tilestep]);
        loq_sn_het += (int)(tilepath.loq_info[0][tilestep].size());
        ez.loq_info_sn_het.push_back(loq_sn_het);
        for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
          ez.loq_info_noc_het.push_back(tilepath.loq_info[0][tilestep][j]);
          ez.loq_info_noc_het.push_back(tilepath.loq_info[0][tilestep][j+1]);
        }

        ez.loq_info_variant_het.push_back(tilepath.allele[1][tilestep]);
        loq_sn_het += (int)(tilepath.loq_info[1][tilestep].size());
        ez.loq_info_sn_het.push_back(loq_sn_het);
        for (j=0; j<tilepath.loq_info[1][tilestep].size(); j+=2) {
          ez.loq_info_noc_het.push_back(tilepath.loq_info[1][tilestep][j]);
          ez.loq_info_noc_het.push_back(tilepath.loq_info[1][tilestep][j+1]);
        }

      }



      ez.loq_info_pos.push_back(tilestep);
      ez.loq_info_variant.push_back(tilepath.allele[0][tilestep]);
      loq_sn += (int)(tilepath.loq_info[0][tilestep].size());
      ez.loq_info_sn.push_back(loq_sn);
      for (j=0; j<tilepath.loq_info[0][tilestep].size(); j+=2) {
        ez.loq_info_noc.push_back(tilepath.loq_info[0][tilestep][j]);
        ez.loq_info_noc.push_back(tilepath.loq_info[0][tilestep][j+1]);
      }

      ez.loq_info_variant.push_back(tilepath.allele[1][tilestep]);
      loq_sn += (int)(tilepath.loq_info[1][tilestep].size());
      ez.loq_info_sn.push_back(loq_sn);
      for (j=0; j<tilepath.loq_info[1][tilestep].size(); j+=2) {
        ez.loq_info_noc.push_back(tilepath.loq_info[1][tilestep][j]);
        ez.loq_info_noc.push_back(tilepath.loq_info[1][tilestep][j+1]);
      }

    }
    tilestep++;
  }
  ez.loq_bv.push_back(u8);

  //
  // -----------------

  // hiq

  // -----------------
  //
  // span_bv
  //

  n_q = n / 8;
  n_r = n % 8;

  tilestep=0;
  for (ii=0; ii<n_q; ii++) {

    u8=0;
    for (i=0; i<8; i++) {

      if (loc_debug) {
        printf("span info tilestep %i, anchor_flag[%i] %i, span_flag[%i] %i\n",
            tilestep,
            tilestep, anchor_flag[tilestep],
            tilestep, span_flag[tilestep]);
      }


      if (anchor_flag[tilestep] || (span_flag[tilestep] != tilestep)) {
        u8 |= (1<<i);
      }
      tilestep++;
    }


    if (loc_debug) {
      printf("pushing span %02x\n", (int)u8);
    }

    ez.span_bv.push_back(u8);
  }

  u8=0;
  for (i=0; i<n_r; i++) {


    if (loc_debug) {
      printf("span info (r) tilestep %i, anchor_flag[%i] %i, span_flag[%i] %i\n",
          tilestep,
          tilestep, anchor_flag[tilestep],
          tilestep, span_flag[tilestep]);
    }


    if (anchor_flag[tilestep] || (span_flag[tilestep] != tilestep)) {

      if (loc_debug) {
        printf("span (r) tilestep %i, anchor_flag[%i] %i, span_flag[%i] %i\n",
            tilestep,
            tilestep, anchor_flag[tilestep],
            tilestep, span_flag[tilestep]);
      }

      u8 |= (1<<i);
    }
    tilestep++;
  }
  ez.span_bv.push_back(u8);

  //
  // -----------------

  // notes:
  // loq_flag now had actual nocall info (spans knots instead of individual tiles)
  // span_bv holds flag to indicate spanning

  n_q = n / CACHE_N;
  n_r = n % CACHE_N;

  if (loc_debug) {
    printf("n_q %i, n_r %i\n", n_q, n_r);
  }

  int ROUGH_COUNT = 0;
  int HEXIT_OVF_COUNT = 0;
  int CACHE_OVF_COUNT = 0;
  int TILEM_OVF_COUNT = 0;
  int SPAN_OVF_COUNT = 0;
  int FILL_OVF_COUNT = 0;
  int t_ovf_count = 0;
  int t_rough_count = 0;


  block_start = 0;
  tilestep=0;
  for (ii=0; ii<n_q; ii++) {

    t_ovf_count=0;
    t_rough_count=0;

    if (loc_debug) {
      printf("--------- block_start: %i\n", block_start);
    }

    u32_canon = 0;
    hexit_vec.clear();
    knot_vec.clear();

    // fill non-anchor spanning tiles with canon bit
    // set.
    //
    for (i=0; i<block_start; i++) { u32_canon |= (1u)<<i; }

    for (i=block_start; i<CACHE_N; i++) {

      block_start=0;

      // If it's loq or spanning, skip (0 bit val)
      //
      /*
      if (loq_flag[tilestep] || (span_flag[tilestep] != tilestep)) {

        if (loc_debug) {
          printf("step %i (%x): loq %i, span %i (loq/span skip)\n", tilestep, tilestep, loq_flag[tilestep], span_flag[tilestep]);
        }

        tilestep++;
        continue;
      }
      */

      // If it's loq, don't set canon bit and move on...
      //
      if (loq_flag[tilestep]) {

        if (loc_debug) {
          printf("step %i (%x): loq %i (loq skip)\n", tilestep, tilestep, loq_flag[tilestep]);
        }

        tilestep++;
        continue;
      }

      // If it's hiq and non-anchor spanning, *set* canon bit.
      // To indicate an anchor step, the span flag is set and the
      // canon bit is not set.  To differentiate anchor spanning tiles from
      // non-anchor spanning tiles, the canon bit is used.
      //
      if (span_flag[tilestep] != tilestep) {

        u32_canon |= (1u)<<i;

        if (loc_debug) {
          printf("step %i (%x): non-anchor spanning (canon bit set)\n", tilestep, tilestep);
        }

        tilestep++;
        continue;
      }

      // If it's hiq and canonical (non-spannig), set bit 1
      //
      if ((!anchor_flag[tilestep]) &&
          (tilepath.allele[0][tilestep]==0) &&
          (tilepath.allele[1][tilestep]==0) ) {

        u32_canon |= (1u)<<i;

        if (loc_debug) {
          printf("step %i (%x): canon\n", tilestep, tilestep);
        }

        tilestep++;
        continue;
      }

      // Finally, if it's hiq, non-canonical and non-anchor spanning, add it to
      // the hexit array for later processing
      //
      if (anchor_flag[tilestep]) {
        k = anchor_flag[tilestep];

        mk_tilemap_key(str, tilepath, tilestep, k);
        tilemap_it = tilemap.find(str);

        //DEBUG
        //fprintf(stderr, "key: %s\n", str.c_str()); fflush(stderr);

        int cf = 0;
        if (tilemap_it != tilemap.end()) {
          // found in tilemap

          if (tilemap_it->second<15) {
            cf = 1;

            hexit_vec.push_back( tilemap_it->second );

            // Add it to overflow array if we've surpassed the
            // space alotted for hexit values.
            //
            if (hexit_vec.size()>=8) { cf=0; }

          } else {

            // found in tilemap but hexit overflow
            //
            hexit_vec.push_back( -3 );


            HEXIT_OVF_COUNT++;
            if (tilemap_it->second<255) {
              ROUGH_COUNT++;
              t_rough_count++;
            }

          }
        } else {

          // tilemap overflow
          //
          hexit_vec.push_back(-4);

          TILEM_OVF_COUNT++;
        }

        if (loc_debug) {
          printf("step %i (%x): hexit %i (len(hexit):%i) anchor ff %i\n", tilestep, tilestep, hexit_vec[ hexit_vec.size()-1 ], (int)hexit_vec.size(), k);
        }

        // cannot put in cache vec, put it in ovf
        //
        if (cf==0) {
          for (j=0; j<k; j++) {
            ez.ovf_vec.push_back((int16_t)(tilestep+j));
            ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep+j]);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep+j]);
            //tilestep++;

            t_ovf_count++;
          }

          SPAN_OVF_COUNT += k;
        }


        if ((i+k) > 32) {

          // non-anchor spanning tiles are set to canon to differentiate
          // anchor spanning from non-anchor spanning
          //
          for (j=i+1; j<32; j++) {
            u32_canon |= (1u)<<j;
          }

          // skip ahead past the spanning positions for this spanning
          // tile
          //
          block_start = (i+k)%32;
        }
        else {

          // non-anchor spanning tiles are set to canon to differentiate
          // anchor spanning from non-anchor spanning
          //
          for (j=i+1; j<(i+k); j++) {
            u32_canon |= (1u)<<j;
          }

          block_start = 0;
        }

        tilestep += k;

        //!!!
        i += k-1;


        continue;
      }

      else {
        sprintf(buf, "%x:%x", tilepath.allele[0][tilestep], tilepath.allele[1][tilestep]);
        str.clear();
        str = buf;

        tilemap_it = tilemap.find(str);

        //DEBUG
        //fprintf(stderr, "key+: %s\n", str.c_str());

        if (tilemap_it != tilemap.end()) {
          if (tilemap_it->second<15) {
            hexit_vec.push_back( tilemap_it->second );

            // add to overflow here if it surpasses the space
            // we have in the hexit cache area
            //
            if (hexit_vec.size()>=8) {
              ez.ovf_vec.push_back((int16_t)tilestep);
              ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep]);
              ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep]);
            }

          } else {
            hexit_vec.push_back(-1);

            ez.ovf_vec.push_back((int16_t)tilestep);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep]);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep]);

            t_ovf_count++;

            HEXIT_OVF_COUNT++;
            if (tilemap_it->second<255) {
              ROUGH_COUNT++;
              t_rough_count++;
            }

          }
        } else {
          hexit_vec.push_back(-2);

          ez.ovf_vec.push_back((int16_t)tilestep);
          ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep]);
          ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep]);

          t_ovf_count++;

          TILEM_OVF_COUNT++;
        }

      }

      if (loc_debug) {
        printf("step %i (%x): hexit %i (len(hexit):%i)\n", tilestep, tilestep, hexit_vec[ hexit_vec.size()-1 ], (int)hexit_vec.size());
      }


      tilestep++;

    }


    if (loc_debug) {
      printf("u32_canon %08x\n", u32_canon);
    }

    //DEBUG
    if (loc_debug) {
      printf("...\n");
      for (k=0; k<hexit_vec.size(); k++) {
        printf("hexit[%i]: %i\n", k, hexit_vec[k]);
      }
      printf("\n");
    }

    u64 = ((uint64_t)u32_canon) << 32;

    u32 = 0;
    int mm = (int)(hexit_vec.size());
    if (mm>8) {
      mm = 8;

      CACHE_OVF_COUNT += (hexit_vec.size()-8);
    }

    else {

      int d = 8 - hexit_vec.size();
      if (d<t_rough_count) { d = t_rough_count; }
      FILL_OVF_COUNT += t_rough_count;

    }

    if (loc_debug) { printf("mm %i\n", mm); }

    uint32_t t32;
    for (i=0; i<mm; i++) {

      if (hexit_vec[i] < 0) {
        //t32 = 0xf;
        //t32 <<= 4*i;
        //u32 |= t32;
        u32 |= (((uint32_t)0xf)<<(4*i));

        if (loc_debug) { printf("hexit ovf: %016x\n", u32); }
      }
      else {
        //t32 = (((uint32_t)hexit_vec[i])&0xf);
        //t32 <<= 4*i;
        //u32 |= t32;
        u32 |= (((uint32_t)hexit_vec[i])&0xf)<<(4*i);

        if (loc_debug) { printf("hexit add: %016x\n", u32); }
      }
    }
    u64 |= (uint64_t)u32;

    ez.cache.push_back(u64);
  }


  if (loc_debug) { printf(">>>>>> block_start: %i\n", block_start); }



  // PROCESS REMAINDER
  //

  u32_canon = 0;

  if (n_r > 0) {

    hexit_vec.clear();
    knot_vec.clear();

    // fill non-anchor spanning tiles with canon bit
    // set.
    //
    for (i=0; i<block_start; i++) {


      if (loc_debug) {
        printf("  span canon bit (tilestep %i, i:%i, block_start:%i)\n", tilestep-block_start+i, i, block_start);
      }

      u32_canon |= (1u)<<i;
    }

    for (i=block_start; i<n_r; i++) {

      if (loc_debug) {
        printf("  remain i%i (tilestep %i)\n", i, tilestep);
      }

      // If it's loq or spanning, skip (0 bit val)
      //
      //if (loq_flag[tilestep] ||
      //    (span_flag[tilestep] != tilestep)) {
      //  tilestep++;
      //  continue;
      //}

      // If it's loq, don't set canon bit and move on...
      //
      if (loq_flag[tilestep]) {

        if (loc_debug) {
          printf("step %i (%x): loq %i (loq skip)\n", tilestep, tilestep, loq_flag[tilestep]);
        }

        tilestep++;
        continue;
      }

      // If it's hiq and non-anchor spanning, *set* canon bit.
      // To indicate an anchor step, the span flag is set and the
      // canon bit is not set.  To differentiate anchor spanning tiles from
      // non-anchor spanning tiles, the canon bit is used.
      //
      if (span_flag[tilestep] != tilestep) {

        u32_canon |= (1u)<<i;

        if (loc_debug) {
          printf("step %i (%x): non-anchor spanning (canon bit set)\n", tilestep, tilestep);
        }

        tilestep++;
        continue;
      }


      // If it's hiq and canonical, set bit 1
      //
      //if ((tilepath.allele[0][tilestep]==0) && (tilepath.allele[1][tilestep]==0)) {
      if ((!anchor_flag[tilestep]) &&
          (tilepath.allele[0][tilestep]==0) &&
          (tilepath.allele[1][tilestep]==0)) {
        u32_canon |= (1u)<<i;

        if (loc_debug) {
          printf("step %i (%x): hiq canonical, set canon bit\n", tilestep, tilestep);
        }

        tilestep++;
        continue;
      }

      // Finally, if it's hiq, non-canonical and non-anchor spanning, add it to
      // the hexit array for later processing
      //
      if (anchor_flag[tilestep]) {
        k = anchor_flag[tilestep];

        mk_tilemap_key(str, tilepath, tilestep, k);
        tilemap_it = tilemap.find(str);

        if (loc_debug) {
          printf("  tilestep %i, key %s\n", tilestep, str.c_str());
        }

        int cf = 0;
        if (tilemap_it != tilemap.end()) {
          // found in tilemap

          if (tilemap_it->second<15) {
            cf = 1;
            hexit_vec.push_back( tilemap_it->second );

            if (loc_debug) {
              printf("    hexit_vec: %i (len %i)\n", tilemap_it->second, (int)hexit_vec.size());
            }

            if (hexit_vec.size()>=8) { cf = 0; }

            if (loc_debug) {
              printf("    (cf %i)\n", cf);
            }

          } else {
            hexit_vec.push_back( -3 );

            if (loc_debug) {
              printf("    hexit_vec: %i (len %i)\n", -3, (int)hexit_vec.size());
            }


          }
        } else {
          hexit_vec.push_back(-4);

          if (loc_debug) {
            printf("    hexit_vec: %i (len %i)\n", -4, (int)hexit_vec.size());
          }

        }

        // cannot put in cache vec, put it in ovf
        //
        if (cf==0) {
          for (j=0; j<k; j++) {
            ez.ovf_vec.push_back((int16_t)tilestep+j);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep+j]);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep+j]);
            //tilestep++;

            if (loc_debug) {
              printf("    ovf: tilestep %i -> %i, %i\n",
                  tilestep+j,
                  tilepath.allele[0][tilestep+j],
                  tilepath.allele[1][tilestep+j] );
            }

          }

        } else {
          //tilestep+=k;
        }

        if ((i+k) > 32) {

          for (j=i+1; j<n_r; j++) {
            u32_canon |= (1u)<<j;
          }

          block_start = (i+k)%32;
        }
        else {

          for (j=i+1; j<(i+k); j++) {
            u32_canon |= (1u)<<j;
          }


          block_start = 0;
        }

        tilestep += k;

        //!!!
        i += k-1;

        continue;
      }

      else {
        sprintf(buf, "%x:%x", tilepath.allele[0][tilestep], tilepath.allele[1][tilestep]);
        str = buf;

        if (loc_debug) {
          printf("  looking up %s\n", str.c_str());
        }

        tilemap_it = tilemap.find(str);
        if (tilemap_it != tilemap.end()) {
          if (tilemap_it->second<15) {

            hexit_vec.push_back( tilemap_it->second );

            if (loc_debug) {
              printf(" hexit++ tilestep %i, tilemapval %i\n", tilestep, tilemap_it->second);
            }

            if (hexit_vec.size()>=8) {

              if (loc_debug) {
                printf(" cache ovf (%i %i %i)\n",
                    tilestep,
                    tilepath.allele[0][tilestep],
                    tilepath.allele[1][tilestep] );
              }

              ez.ovf_vec.push_back((int16_t)tilestep);
              ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep]);
              ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep]);
            }

          } else {
            hexit_vec.push_back(-1);

            if (loc_debug) {
              printf(" val ovf (%i %i %i)\n",
                  tilestep,
                  tilepath.allele[0][tilestep],
                  tilepath.allele[1][tilestep] );
            }


            ez.ovf_vec.push_back((int16_t)tilestep);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep]);
            ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep]);

          }
        } else {
          hexit_vec.push_back(-2);

          if (loc_debug) {
            printf(" tilemap ovf (%i %i %i)\n",
                tilestep,
                tilepath.allele[0][tilestep],
                tilepath.allele[1][tilestep] );
          }



          ez.ovf_vec.push_back((int16_t)tilestep);
          ez.ovf_vec.push_back((int16_t)tilepath.allele[0][tilestep]);
          ez.ovf_vec.push_back((int16_t)tilepath.allele[1][tilestep]);
        }

      }

      tilestep++;

    }

    u64 = ((uint64_t)u32_canon) << 32;

    u32 = 0;
    int mm = (int)(hexit_vec.size());
    if (mm>8) { mm = 8; }

    if (loc_debug) { printf("mm %i\n", mm); }

    for (i=0; i<mm; i++) {



      if (hexit_vec[i] < 0) {
        u32 |= (((uint32_t)0xf)<<(4*i));

        if (loc_debug) { printf("hexit ovf: %016x\n", u32); }
      }
      else {
        u32 |= (((uint32_t)hexit_vec[i])&0xf)<<(4*i);

        if (loc_debug) { printf("hexit add: %016x\n", u32); }
      }
    }
    u64 |= (uint64_t)u32;

    ez.cache.push_back(u64);

    if (loc_debug) {
      printf("pushing cache %08llx\n", (long long unsigned int)u64);
    }

  }


  //DEBUG
  if (loc_debug) {
    printf("ROUGH_COUNT: %i\n", ROUGH_COUNT);
    printf("TILEM_OVF_COUNT: %i\n", TILEM_OVF_COUNT);
    printf("HEXIT_OVF_COUNT: %i\n", HEXIT_OVF_COUNT);
    printf("CACHE_OVF_COUNT: %i\n", CACHE_OVF_COUNT);
    printf("SPAN_OVF_COUNT: %i\n", SPAN_OVF_COUNT);
    printf("FILL_OVF_COUNT: %i\n", FILL_OVF_COUNT);
  }

}

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
int old_load_tilemap(const char *fn_tm, std::map< std::string, int > &tilemap) {
  int i, j, k, n;
  int pos;
  char buf[1024];
  FILE *fp;

  std::map< std::string, int >::iterator ent;
  std::string key;

  fp = fopen(fn_tm, "r");
  if (!fp) { perror(fn_tm); exit(1); }

  pos = 0;
  while (fgets(buf, 1023, fp)) {
    if (feof(fp)) { break; }

    n = strlen(buf);
    if (n<2) { continue; }

    buf[n-1] = '\0';

    key = buf;

    tilemap[key] = pos;
    pos++;

  }

  fclose(fp);

}

void print_tilepath_vec(tilepath_vec_t &tv) {
  int i, j, k;
  printf("%s\n", tv.name.c_str());
  for (i=0; i<2; i++) {
    for (j=0; j<tv.allele[i].size(); j++) {
      printf(" %i [%i]", tv.allele[i][j], tv.loq_flag[i][j]);

      printf("{");
      for (k=0; k<tv.loq_info[i][j].size(); k++) {
        printf(" %i", tv.loq_info[i][j][k]);
      }
      printf("}");

    }
    printf("\n");
  }
}

void print_bgf(tilepath_vec_t &tv) {
  int a=0;
  int i, j, k;

  for (a=0; a<2; a++) {
    printf("[");
    for (i=0; i<tv.allele[a].size(); i++) {
      printf(" %i", tv.allele[a][i]);
    }
    printf("]\n");
  }

  for (a=0; a<2; a++) {
    printf("[");
    for (i=0; i<tv.loq_info[a].size(); i++) {
      printf("[");
      for (j=0; j<tv.loq_info[a][i].size(); j++) {
        printf(" %i", tv.loq_info[a][i][j]);
      }
      printf(" ]");
    }
    printf("]\n");
  }

}

void print_tilepath_vecs(std::vector<tilepath_vec_t> &tv) {
  size_t n, tilepath_n;
  int i, j, pos;
  int max_val;
  int len;

  int print_header = 1;
  int hotpos;

  n = tv.size();
  tilepath_n = tv[0].allele[0].size();

  for (pos=0; pos<tilepath_n; pos++) {

    max_val = 0;
    for (i=0; i<n; i++) {

      if ((tv[i].loq_flag[0][pos]==0) &&
          (max_val < tv[i].allele[0][pos])) {
        max_val = tv[i].allele[0][pos];
      }

      if ((tv[i].loq_flag[1][pos]==0) &&
          (max_val < tv[i].allele[1][pos])) {
        max_val = tv[i].allele[1][pos];
      }

    }

    len = (int)(max_val)+1;

    for (hotpos=0; hotpos<len; hotpos++) {
      for (i=0; i<n; i++) {

        if (i>0) { printf(" "); }
        else if (print_header) {
          //printf("pos%03x.%03x.u | ", pos, hotpos);
          if (TILEPATH>=0) {
            printf("%04x.%03x(%03x)u ", TILEPATH, pos, hotpos);
          } else {
            printf("pos%03x(%03x)u ", pos, hotpos);
          }
        }

        if (tv[i].loq_flag[0][pos]==0) {
          printf("%i", (tv[i].allele[0][pos] == hotpos) ? 1 : 0);
        } else {
          printf("NaN");
        }

      }

      printf("\n");

    }

    for (hotpos=0; hotpos<len; hotpos++) {

      for (i=0; i<n; i++) {

        if (i>0) { printf(" "); }
        else if (print_header) {
          //printf("pos%03x.%03x.v | ", pos, hotpos);
          if (TILEPATH>=0) {
            printf("%04x.%03x(%03x)v ", TILEPATH, pos, hotpos);
          } else {
            printf("pos%03x(%03x)v ", pos, hotpos);
          }
        }

        if (tv[i].loq_flag[1][pos]==0) {
          printf("%i", (tv[i].allele[1][pos] == hotpos) ? 1 : 0);
        } else {
          printf("NaN");
        }

      }

      printf("\n");
    }

    //printf("\n");
  }

}

void print_tilemap(std::map< std::string, int > &tilemap) {
  std::map< std::string, int >::iterator ent;
  std::string s;
  int val;

  for (ent = tilemap.begin(); ent != tilemap.end(); ++ent ) {
    s = ent->first;
    val = ent->second;

    printf(">> %i %s\n", val, s.c_str());
  }

}
