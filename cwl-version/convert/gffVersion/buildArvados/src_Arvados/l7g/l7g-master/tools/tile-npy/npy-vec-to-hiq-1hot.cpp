/* create 'flat' numpy hiq tile vector arrays and it's info file.
 * create 1hot numpy hiq tile vector arrays and it's info file.
 *
 * output:
 *
    [out_dir]/hiq
    [out_dir]/hiq-info

    [out_dir]/hiq-1hot
    [out_dir]/hiq-1hot-info
 *
 *
 * format of 'info' is:
 *
 *    (tilepath << 20) + (tilestep*2) + (allele)
 *
 * for example, info value 891342211 (0x3520cd83) is tilepath 850 (0x352),
 * tile step 26305 (0x66c1) allele 1.
 *
 */

// to compile:
//
//    g++  -O3 -static npy-vec-to-hiq-1hot.cpp  -o npy-vec-to-hiq-1hot -L../../lib/cnpy -lcnpy -I../../lib/cnpy
//
// to run:
//
//   LD_LIBRARY_PATH=$cnpylibdir ./npy-vec-to-hiq-1hot names.npy inp_vec_npy_dir/ [out_dir/]
//

#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <string.h>

#include <cstdlib>
#include <iostream>
#include <map>
#include <string>

#include <vector>
#include <string>
#include <iterator>
#include <complex>

#include "cnpy.h"

int append_1hot(std::vector<int> &ohv, int val, int max_val) {
  int i, j, k;
  for (i=0; i<max_val; i++) {
    if (i==val) { ohv.push_back(1); }
    else { ohv.push_back(0); }
  }
}

int inv_1hot(std::vector<int> &ohv, int s, int n) {
  int i, j, k;
  int val = -1;
  int pop_count = 0;

  for (i=s; i<(s+n); i++) {
    if (ohv[i]) {
      val = i-s;
      pop_count++;
    }
  }

  if (pop_count>1) { return -2; }
  return val;
}

int check_1hot(std::vector< std::vector<int> > &hiq_ilv, std::vector<int> &hiq_pos,
    std::vector< std::vector<int> > &oh_hiq, std::vector<int> &oh_hiq_info) {

  int i, j, k;
  int n_ds, ds, n_vec, pos, n_oh_vec, start_oh_v_pos;
  int t, m;

  int v_pos, oh_v_pos, prev_oh_v_pos;

  int verbose = 1;

  if (verbose) { printf("  check_1hot: cp0\n"); fflush(stdout); }

  if (hiq_ilv.size() != oh_hiq.size()) { return -1; }

  if(verbose) { printf("  check_1hot: cp2\n"); fflush(stdout); }

  for (ds=1; ds<hiq_ilv.size(); ds++) {
    if (hiq_ilv[ds-1].size() != hiq_ilv[ds].size()) { return -2; }
  }

  if (verbose) { printf("  check_1hot: cp3\n"); fflush(stdout); }

  if (hiq_ilv[0].size() != hiq_pos.size()) { return -3; }

  if (verbose) { printf("  check_1hot: cp4\n"); fflush(stdout); }

  for (ds=1; ds<oh_hiq.size(); ds++) {
    if (oh_hiq[ds-1].size() != oh_hiq[ds].size()) { return -4; }
  }

  if (verbose) { printf("  check_1hot: cp5\n"); fflush(stdout); }

  if (oh_hiq[0].size() != oh_hiq_info.size()) {  return -5; }

  if (verbose) { printf("  check_1hot: cp6\n"); fflush(stdout); }


  n_ds = (int)(hiq_ilv.size());
  n_vec = (int)(hiq_ilv[0].size());
  n_oh_vec = (int)(oh_hiq[0].size());

  oh_v_pos = 0;
  for (v_pos=0; v_pos<n_vec; v_pos++) {
    if (oh_v_pos>=n_oh_vec) { break; }

    if (oh_hiq_info[oh_v_pos] != hiq_pos[v_pos]) { return -6; }
    while ((oh_v_pos < n_oh_vec) && (oh_hiq_info[oh_v_pos] == hiq_pos[v_pos])) {
      oh_v_pos++;
    }
  }

  if (verbose) { printf("  check_1hot: cp7\n"); fflush(stdout); }

  if (verbose) {
    printf(" oh_v_pos %i, n_oh_vec %i, v_pos %i, n_vec %i\n",
        oh_v_pos, n_oh_vec, v_pos, n_vec);
    fflush(stdout);
    if (v_pos<n_vec) {
      for (i=v_pos; i<n_vec; i++) { printf(" (%i,%i)", i, hiq_pos[i]); }
      printf("\n");
      fflush(stdout);
    }
  }

  if (v_pos!=n_vec) { return -7; }

  if (verbose) {
    printf("  check_1hot: cp8\n"); fflush(stdout);
  }

  if (oh_v_pos != n_oh_vec) { return -8; }

  if (verbose) {
    printf("  check_1hot: cp9\n"); fflush(stdout);
  }

  for (ds=0; ds<n_ds; ds++) {

    oh_v_pos = 0;
    for (v_pos=0; v_pos<n_vec; v_pos++) {

      start_oh_v_pos = oh_v_pos;
      while ((oh_v_pos < n_oh_vec) && (oh_hiq_info[oh_v_pos] == hiq_pos[v_pos])) {
        oh_v_pos++;
      }

      k = inv_1hot( oh_hiq[ds], start_oh_v_pos, oh_v_pos-start_oh_v_pos);
      if (k<-1) {

        if (verbose) {
        printf("INVALID VAL: got val %i expecting %i @ vpos %i, oh_vpos %i+%i, ds %i\n",
            k, hiq_ilv[ds][v_pos],
            v_pos, start_oh_v_pos, oh_v_pos - start_oh_v_pos,
            ds);
        }

        return -10;
      }

      if (k!=hiq_ilv[ds][v_pos]) {

        if (verbose) {
        printf("MISMATCH: got val %i expecting %i @ vpos %i, oh_vpos %i+%i, ds %i\n",
            k, hiq_ilv[ds][v_pos],
            v_pos, start_oh_v_pos, oh_v_pos - start_oh_v_pos,
            ds);
        }

        return -9;
      }

    }

  }

  if (verbose) { printf("  check_1hot: cp10\n"); fflush(stdout); }

  return 0;
}

int write_1hot(std::vector< std::vector<int> > &hiq_ilv, std::vector<int> &hiq_pos, std::string &odir) {
  int i, j, k, ds, vpos, n_dataset, n_vec;
  int *tvec;
  std::vector< std::vector<int> > oh_hiq;
  std::vector<int> oh_hiq_info;
  int oh_max;
  int shape[2];
  size_t oh_width, idx;

  std::string ofn, ofn_info;

  ofn = odir;
  ofn += "/hiq-1hot";
  ofn_info = odir;
  ofn_info += "/hiq-1hot-info";

  n_dataset = hiq_ilv.size();
  n_vec = hiq_ilv[0].size();


  // first calc max
  //
  int n_tot=0;
  std::vector<int> max_val;

  for (vpos=0; vpos<n_vec; vpos+=2) {
    oh_max=0;
    for (ds=0; ds<n_dataset; ds++) {
      if (oh_max<hiq_ilv[ds][vpos]) { oh_max = hiq_ilv[ds][vpos]; }
      if (oh_max<hiq_ilv[ds][vpos+1]) { oh_max = hiq_ilv[ds][vpos+1]; }
    }
    max_val.push_back(oh_max);
    n_tot += oh_max+1;
  }

  for (ds=0; ds<n_dataset; ds++) {
    std::vector<int> v;
    oh_hiq.push_back(v);
  }

  for (ds=0; ds<n_dataset; ds++) {
    for (vpos=0; vpos<n_vec; vpos++) {
      oh_max=max_val[vpos/2];
      append_1hot(oh_hiq[ds], hiq_ilv[ds][vpos], oh_max+1);
    }
  }

  for (vpos=0; vpos<n_vec; vpos++) {
    for (i=0; i<=max_val[vpos/2]; i++) {
      oh_hiq_info.push_back(hiq_pos[vpos]);
    }
  }

  printf("calling check_1hot...\n"); fflush(stdout);

  k = check_1hot(hiq_ilv, hiq_pos, oh_hiq, oh_hiq_info);
  printf("GOT: %i\n", k);
  if (k<0) { exit(k); }


  printf("...sanity...\n");
  for (ds=1; ds<n_dataset; ds++) {
    if (oh_hiq[ds-1].size() != oh_hiq[ds].size()) {
      printf("sanity error on ds %i\n", ds);
      exit(-1);
    }
  }
  printf("...ok\n");
  fflush(stdout);

  printf("... %i %i, %i\n",
      (int)(oh_hiq.size()),
      (int)(oh_hiq[0].size()),
      (int)(oh_hiq_info.size()));
  fflush(stdout);

  oh_width = (size_t)oh_hiq[0].size();

  double *dvec = (double *)malloc(sizeof(double)*(oh_hiq.size()*oh_hiq[0].size()));
  for (ds=0; ds<n_dataset; ds++)  {
    for (vpos=0; vpos<oh_hiq[0].size(); vpos++) {
      idx = oh_width;
      idx *= (size_t)ds;
      idx += vpos;
      //dvec[n_vec*ds + vpos] = (double)(oh_hiq[ds][vpos]);
      dvec[idx] = (double)(oh_hiq[ds][vpos]);
    }
  }

  printf("cp0\n"); fflush(stdout);

  shape[0] = (int)(oh_hiq.size());
  shape[1] = (int)(oh_hiq[0].size());
  cnpy::npy_save(ofn, dvec, (const unsigned int *)shape, 2, "w");
  free(dvec);

  printf("cp1\n"); fflush(stdout);

  shape[0] = oh_hiq[0].size();
  cnpy::npy_save(ofn_info, &(oh_hiq_info[0]), (const unsigned int *)shape, 1, "w");



  printf("n_tot %i, n_dataset %i (%i)\n", n_tot, n_dataset, n_tot*n_dataset);

}

int load_list(const char *fn, std::vector< std::string > &name_vec) {
  char buf[1001];
  FILE *fp;
  int i, j, k, n;
  char *p;

  std::string s;

  fp = fopen(fn, "r");
  while (!feof(fp)) {
    p = fgets(buf, 1000, fp);

    if (!p) { continue; }

    n = strlen(buf);

    if (n==0) { continue; }
    if (buf[0]=='\n') { continue; }
    if (buf[n-1]=='\n') { buf[n-1] = '\0'; }



    s = buf;
    name_vec.push_back(s);
  }

  fclose(fp);

  return 0;

  /*
  for (i=0; i<name_vec.size(); i++) { printf("%s\n", name_vec[i].c_str()); }
  printf("EXITING");
  exit(0);
  */

}

int load_all(std::string &npy_name_list, std::string &npy_vec_idir, std::string &odir) {
  int i, j, k;
  int n_dataset,n_vec, n_out_dataset;
  std::string s;
  char buf[1024];
  cnpy::NpyArray raw, names;

  int cur_ds;

  std::vector< std::vector<int> > hiq_ilv;
  std::vector<int> hiq_pos;

  int *tvec, *t_info_vec;
  int tilepath, vpos, ds, is_loq=0;
  int hiq_count=0;
  int shape[2];

  std::string ofn, ofn_info;

  ofn = odir;
  ofn += "/hiq";
  ofn_info = odir;
  ofn_info += "/hiq-info";

  names = cnpy::npy_load(npy_name_list.c_str());

  for (i=0; i<names.shape[0]; i++) {
    std::vector<int> v;
    hiq_ilv.push_back(v);
  }

  for (tilepath=0; tilepath<=862; tilepath++) {
    sprintf(buf, "%s/%03x", npy_vec_idir.c_str(), tilepath);

    raw = cnpy::npy_load(buf);
    n_dataset = (int)(raw.shape[0]);
    n_vec = (int)(raw.shape[1]);
    tvec = reinterpret_cast<int *>(raw.data);

    hiq_count=0;

    for (vpos=0; vpos<n_vec; vpos+=2) {
      is_loq=0;
      for (cur_ds=0, ds=0; ds<n_dataset; ds++) {

        if ( tvec[n_vec*ds + vpos] == -2 ) {
          is_loq=1;
          break;
        }
        if ( tvec[n_vec*ds + vpos+1] == -2 ) {
          is_loq=1;
          break;
        }
      }
      if (is_loq) { continue; }

      hiq_pos.push_back( (tilepath << 20) + vpos );
      hiq_pos.push_back( (tilepath << 20) + vpos + 1);

      for (cur_ds=0, ds=0; ds<n_dataset; ds++) {

        hiq_ilv[cur_ds].push_back( tvec[n_vec*ds + vpos] );
        hiq_ilv[cur_ds].push_back( tvec[n_vec*ds + vpos + 1] );

        cur_ds++;
      }

      hiq_count++;
    }


    printf("[%03x]: got %i %i, hiq_count %i\n", tilepath, n_dataset, n_vec, hiq_count);

    raw.destruct();
  }

  n_dataset = (int)(hiq_ilv.size());
  n_vec = (int)(hiq_ilv[0].size());


  // save integer tiling vectors
  //

  printf("cp0.5\n"); fflush(stdout);
  printf("n_dataset: %i, n_vec: %i\n", n_dataset, n_vec);
  printf("%i %i %i\n",
      (int)(hiq_ilv[0].size()) ,
      (int)(hiq_ilv[1].size()) ,
      (int)(hiq_ilv[2].size()) ); fflush(stdout);

  tvec = (int *)malloc(sizeof(int)* n_vec * n_dataset);
  for (ds=0; ds<n_dataset; ds++)  {
    for (vpos=0; vpos<n_vec; vpos++) {
      tvec[n_vec*ds + vpos] = hiq_ilv[ds][vpos];
    }
  }

  printf("cp0.7\n"); fflush(stdout);
  printf("sanity: %i %i, %i\n",
      (int)(hiq_pos.size()),
      (int)(hiq_ilv.size()),
      (int)(hiq_ilv[0].size()));
  for (i=1; i<hiq_ilv.size(); i++) {
    if (hiq_ilv[i-1].size() != hiq_ilv[i].size()) {
      printf("error! %i (%i != %i)\n", i, (int)(hiq_ilv[i-1].size()), (int)(hiq_ilv[i].size()) );
      exit(-1);
    }
  }

  printf("cp0\n"); fflush(stdout);

  shape[0] = n_dataset;
  shape[1] = n_vec;
  cnpy::npy_save(ofn, tvec, (const unsigned int *)shape, 2, "w");
  free(tvec);

  printf("cp1\n"); fflush(stdout);
  printf("writing shape %i (%i)\n", n_vec, (int)(hiq_pos.size()));
  printf("cp2\n"); fflush(stdout);
  printf("hiq_pos.size() %i, n_vec %i\n", (int)hiq_pos.size(), (int)n_vec);
  printf("  ofn_info: %s\n", ofn_info.c_str());

  shape[0] = n_vec;
  cnpy::npy_save(ofn_info, &(hiq_pos[0]), (const unsigned int *)shape, 1, "w");

  // now calculate 1hot
  //
  write_1hot( hiq_ilv, hiq_pos, odir );

  names.destruct();
}

int main(int argc, char **argv) {

  std::string npy_name_list, npy_vec_idir, odir;

  if (argc<3) {
    printf("provide npy-name-list and npy-vec-tilepath-dir\n");
    exit(-1);
  }

  npy_name_list = argv[1];
  npy_vec_idir = argv[2];
  odir = "./";

  if (argc>=4) {
    odir = argv[3];
  }

  printf("npy_name_list %s\n", npy_name_list.c_str());
  printf("npy_vec_idir %s\n", npy_vec_idir.c_str());
  printf("odir: %s\n", odir.c_str());

  load_all( npy_name_list, npy_vec_idir, odir );
}
