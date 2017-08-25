/* Consolidate individual Lightning numpy arrays broken out by tilepath
 * into a single Lightning tile numpy matrix.
 *
 * to run:
 *
 *   ./npy-consolidate inp-data-vec/[0123]* out.npy
 *
  *
 */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "cnpy.h"

#include <vector>
#include <string>



int main(int argc, char **argv) {
  int i, j, k;
  std::vector< std::string > ifns;
  std::vector< std::vector< int > > ivec;
  std::vector< int > tvec;
  cnpy::NpyArray npyvec;
  std::string s;

  std::vector< int > npy_info;
  std::string ofn, ofn_info;

  int loc_debug=0;

  int *pvec;
  uint32_t u32;
  uint64_t u64;

  size_t sz, idx_sz;

  int shape[2];
  int *out_vec, start_pos=0;
  int idx, ds, vpos;

  int tile_path, tile_step_x2, val;

  ofn = "out-consolidated";

  if (argc<=2) {
    printf("provide input npy list\n");
    printf("\n");
    printf("example:\n\n  ./npy-consolidate inp-data-vec/[0123]* out.npy\n\n");
    exit(-1);
  }

  ifns.clear();


  for (i=1; i<(argc-1); i++) { ifns.push_back(argv[i]); }
  ofn = argv[argc-1];


  shape[0] = 0;
  shape[1] = 0;

  if (loc_debug) {
    for (i=0; i<ifns.size(); i++) {
      printf(" ifn[%i] %s\n", i, ifns[i].c_str());
    }
    printf("ofn: %s\n", ofn.c_str());
  }


  tile_path=0;
  for (i=0; i<ifns.size(); i++) {

    npyvec = cnpy::npy_load(ifns[i].c_str());

    pvec = reinterpret_cast<int *>(npyvec.data);

    if (i==0) {
      shape[0] = npyvec.shape[0];
      for (idx=0; idx<shape[0]; idx++) {
        ivec.push_back(tvec);
      }
    }

    for (ds=0; ds<npyvec.shape[0]; ds++) {

      tile_step_x2=0;
      for (vpos=0; vpos<npyvec.shape[1]; vpos++) {
        ivec[ds].push_back( pvec[ ds*npyvec.shape[1] + vpos ] );

        if (ds==0) {
          val = tile_path<<20;
          val |= tile_step_x2;
          tile_step_x2++;
          npy_info.push_back(val);
        }

      }

    }

    if (shape[0] != npyvec.shape[0]) {
      fprintf(stderr, "ROW SHAPES DO NOT MATCH for idx %i (%s).  exiting\n",
          i, ifns[i].c_str());
      exit(-1);
    }

    shape[0] = npyvec.shape[0];
    shape[1] += npyvec.shape[1];

    npyvec.destruct();

    printf("tile_path: %i %i\n", i, tile_path);
    tile_path++;
  }

  if (loc_debug) {
    printf("shape %i %i\n", shape[0], shape[1]);
    printf("ivec %i %i\n", (int)ivec.size(), (int)ivec[0].size());
    fflush(stdout);

    for (i=1; i<ivec.size(); i++) {
      if (ivec[i-1].size() != ivec[i].size()) {
        fprintf(stderr, "[%i] size (%i) != [%i] (%i)\n",
            i-1, (int)ivec[i-1].size(),
            i, (int)ivec[i].size());
      }
    }
    printf("ok..\n"); fflush(stdout);
  }

  sz = (size_t)shape[0];
  sz *= (size_t)shape[1];

  out_vec = (int *)malloc(sizeof(int)*sz);

  if (loc_debug) {
    printf("%p\n", out_vec);
    printf("%llu\n", (unsigned long long int)sz);
    for (idx_sz=0; idx_sz<sz; idx_sz++) {
      out_vec[idx_sz] = 0;
    }
    printf("...\n"); fflush(stdout);
  }


  start_pos = 0;
  for (ds=0; ds<shape[0]; ds++) {
    for (vpos=0; vpos<shape[1]; vpos++) {
      idx_sz = (size_t)(ds);
      idx_sz *= (size_t)(shape[1]);
      idx_sz += vpos;
      out_vec[idx_sz] = ivec[ds][vpos];
    }
  }

  if (loc_debug) {
    printf("cp\n"); fflush(stdout);
  }

  cnpy::npy_save(ofn, out_vec, (const unsigned int *)shape, 2, "w");

  shape[0] = (int)npy_info.size();
  ofn_info = ofn + "-info";
  cnpy::npy_save(ofn_info, (const unsigned int *)(&(npy_info[0])), (const unsigned int *)shape, 1, "w");

  free(out_vec);
}
