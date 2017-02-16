#include "cgft.hpp"

#define CGFT_VERSION "0.1.0"

#define cleanup_err() do { ret=1; show_help(); goto cgft_cleanup; } while (0);
#define cleanup_fail() do { ret=-2; goto cgft_cleanup; } while (0);
#define cleanup_ok() do { ret=0; show_help(); goto cgft_cleanup; } while (0);
#define cleanup() do { show_help(); goto cgft_cleanup; } while (0);


void cgft_stats(cgf_t *cgf) {

  tilepath_t *tilepath;
  uint64_t u64, *cache;
  uint32_t canon_bit, u32, hiq_mask;
  int n_8, n_32, ntile, s, p, run_sum=0, n;
  unsigned char *loq, *span;
  uint16_t *ovf;

  size_t span_bv_size=0, loq_bv_size=0, cache_size=0, ovf_size=0, loq_size=0;
  size_t sz=0;
  int ovf_count=0, tot_tile=0;

  for (p=0; p<cgf->Path.size(); p++) {
    tilepath = &(cgf->Path[p]);

    ntile = (int)tilepath->NTileStep;
    n_8 = ((ntile+7)/8);
    n_32 = ((ntile+31)/32);

    tot_tile += ntile;

    loq = tilepath->Loq;
    span = tilepath->Span;
    cache = tilepath->Cache;
    ovf = tilepath->Overflow;


    loq_bv_size += (size_t)n_8;
    span_bv_size += (size_t)n_8;
    cache_size += (size_t)(n_32*sizeof(uint64_t));
    ovf_size += (size_t)(tilepath->NOverflow * sizeof(uint16_t));

    ovf_count += (int)(tilepath->NOverflow / 3);

    loq_size += (size_t)tilepath->LoqTileStepHomSize;
    loq_size += (size_t)tilepath->LoqTileVariantHomSize;
    loq_size += (size_t)tilepath->LoqTileNocSumHomSize;
    loq_size += (size_t)tilepath->LoqTileNocStartHomSize;
    loq_size += (size_t)tilepath->LoqTileNocLenHomSize;

    loq_size += (size_t)tilepath->LoqTileStepHetSize;
    loq_size += (size_t)tilepath->LoqTileVariantHetSize;
    loq_size += (size_t)tilepath->LoqTileNocSumHetSize;
    loq_size += (size_t)tilepath->LoqTileNocStartHetSize;
    loq_size += (size_t)tilepath->LoqTileNocLenHetSize;


  }

  float denom = (float)( tot_tile ? tot_tile : 1 );
  float r = (float)ovf_count / denom;

  printf("cgf stats:\n");
  printf("loq_bv: %u\n", (unsigned int)loq_bv_size);
  printf("span_bv: %u\n", (unsigned int)span_bv_size);
  printf("cache: %u\n", (unsigned int)cache_size);
  printf("ovf: %u (%i / %i, %f)\n", (unsigned int)ovf_size, ovf_count, tot_tile, r);
  printf("loq: %u\n", (unsigned int)loq_size);
}


void cgft_run_test(cgf_t *cgf) {
  int i, j, k;
  std::vector< std::vector<int> > tmap[2];
  std::vector<int> tspan;
  std::vector<int> v;

  timeval tv_start, tv_end;

  gettimeofday(&tv_start, NULL);

  v.clear(); v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(1); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(1); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(1); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(1); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(2); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(2); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(0); v.push_back(1);
             v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(1); v.push_back(2); tmap[1].push_back(v);
  tspan.push_back(2);

  v.clear(); v.push_back(1); v.push_back(2); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(1);
             v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(2);

  v.clear(); v.push_back(0); v.push_back(2); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(2); tmap[1].push_back(v);
  tspan.push_back(2);

  v.clear(); v.push_back(0); v.push_back(2); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(2); tmap[1].push_back(v);
  tspan.push_back(2);

  v.clear(); v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(3); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(3); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(1); v.push_back(1);
             v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(2); tmap[1].push_back(v);
  tspan.push_back(2);

  v.clear(); v.push_back(0); v.push_back(2); tmap[0].push_back(v);
  v.clear(); v.push_back(1); v.push_back(1);
             v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(2);

  v.clear(); v.push_back(0); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(4); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(4); v.push_back(1); tmap[0].push_back(v);
  v.clear(); v.push_back(0); v.push_back(1); tmap[1].push_back(v);
  tspan.push_back(1);

  v.clear(); v.push_back(1); v.push_back(2); tmap[0].push_back(v);
  v.clear(); v.push_back(1); v.push_back(2); tmap[1].push_back(v);
  tspan.push_back(2);


  /*
  for (i=0; i<tspan.size(); i++) {
    for (j=0; j<tmap[0][i].size(); j+=2) {
      printf(" %i+%i", tmap[0][i][j], tmap[0][i][j+1]);
    }
    printf(":");
    for (j=0; j<tmap[1][i].size(); j+=2) {
      printf(" %i+%i", tmap[1][i][j], tmap[1][i][j+1]);
    }

    printf(" (%i)\n", tspan[i]);
  }

  printf("\n\n");
  */

  tilepath_t *tilepath;
  uint64_t u64, *cache;
  uint32_t canon_bit, u32, hiq_mask;
  int n_8, n_32, ntile, s, p, run_sum=0, n;
  unsigned char *loq, *span;
  uint16_t *ovf;

  for (p=0; p<cgf->Path.size(); p++) {
    tilepath = &(cgf->Path[p]);

    ntile = (int)tilepath->NTileStep;
    n_8 = ((ntile+7)/8);
    n_32 = ((ntile+31)/32);

    loq = tilepath->Loq;
    span = tilepath->Span;
    cache = tilepath->Cache;
    ovf = tilepath->Overflow;

    for (s=0; s<n_32; s++) {
      hiq_mask = (loq[s*4]) | (loq[s*4+1]<<8) | (loq[s*4+2]<<16) | (loq[s*4+3]<<24);
      hiq_mask = ~hiq_mask;
      u32 = (uint32_t)(cache[s] >> 32);
      hiq_mask &= cache[s];

      run_sum += NumberOfSetBits32(hiq_mask);
    }

    n = (int)tilepath->NOverflow;
    for (i=0; i<n; i+=3) {
      if (ovf[i+1] == OVF16_MAX)  { run_sum++; }
      else                      { run_sum += (int)ovf[i+1]; }

      if (ovf[i+2] == OVF16_MAX)  { run_sum++; }
      else                      { run_sum += (int)ovf[i+2]; }
    }


  }


  gettimeofday(&tv_end, NULL);

  printf("run_sum: %i\n", run_sum);

  float f = (float)(tv_end.tv_usec - tv_start.tv_usec) / 1000.0;

  printf("beg: %i.%i\n", (int)tv_start.tv_sec, (int)tv_start.tv_usec);
  printf("end: %i.%i\n", (int)tv_end.tv_sec, (int)tv_end.tv_usec);
  printf("...: %i (%f)\n", (int)(tv_end.tv_usec - tv_start.tv_usec), f);

}

void show_version() {
  printf("%s\n", CGFT_VERSION);
}

void show_help() {
  printf("CGF Tool.  A tool used to inspect and edit Compact Genome Format (CGF) files.\n");
  printf("Version: %s\n", CGFT_VERSION);
  printf("\n");
  //printf("usage:\n");
  //printf("\n");
  printf("usage: cgft [-H] [-b tilepath] [-e tilepath] [-i ifn] [-o ofn] [-h] [-v] [-V] [ifn]\n");
  printf("\n");
  printf("  [-H|--header]           show header\n");
  printf("  [-C]                    create empty container\n");
  printf("  [-I]                    print basic information about CGF file\n");
  printf("  [-b|--band tilepath]    output band for tilepath\n");
  printf("  [-e|--encode tilepath]  input tilepath band and add it to file, overwriting if it already exists\n");
  printf("  [-i ifn]                input file (CGF)\n");
  printf("  [-o ofn]                output file (CGF)\n");
  printf("  [-A]                    show all tilepaths\n");
  printf("  [-h]                    show help (this screen)\n");
  printf("  [-v]                    show version\n");
  printf("  [-V]                    set verbose level\n");
  printf("  [-Z]                    print \"ez\" structure information\n");
  printf("\n");
}

int main(int argc, char **argv) {
  int i, j, k;
  int opt;
  int ret=0;
  int idx;
  char buf[1024];

  int opt_show_header = 0, opt_show_band=0, opt_encode=0,
      opt_show_help=0, opt_show_version=0, opt_verbose=0,
      opt_delete=0, opt_create_container=0, opt_tilemap=0,
      opt_show_all=0, opt_ez_print=0;
  int opt_run_test=0, opt_info=0;
  char *opt_ifn=NULL, *opt_ofn=NULL, *opt_tilemap_fn=NULL;
  char *opt_band_ifn=NULL;
  FILE *opt_band_ifp = NULL;

  int opt_tilepath=-1;

  std::string tilemap_str;

  FILE *ifp=NULL, *ofp=NULL;

  cgf_t *cgf=NULL;

  while ((opt = getopt(argc, argv, "Hb:e:i:o:CT:hvVAZRI"))!=-1) switch (opt) {
    case 'H': opt_show_header=1; break;
    case 'C': opt_create_container=1; break;
    case 'I': opt_info=1; break;
    case 'T': opt_tilemap=1; opt_tilemap_fn=strdup(optarg); break;
    case 'b': opt_show_band=1; opt_tilepath=atoi(optarg); break;
    case 'e': opt_encode=1; opt_tilepath=atoi(optarg); break;
    case 'd': opt_delete=1; opt_tilepath=atoi(optarg); break;
    case 'i': opt_ifn=strdup(optarg); break;
    case 'o': opt_ofn=strdup(optarg); break;
    case 'h': opt_show_help=1; break;
    case 'A': opt_show_all=1; break;
    case 'v': opt_show_version=1; break;
    case 'V': opt_verbose=1; break;
    case 'Z': opt_ez_print=1; break;
    case 'R': opt_run_test=1; break;
    default: printf("unknown option"); show_help(); cleanup_ok(); break;
  }

  if (argc>optind) {
    if ((argc-optind)>1) { printf("Extra options specified\n"); cleanup_err(); }
    if (opt_ifn) { printf("Input CGF already specified.\n"); cleanup_err(); }
    opt_ifn = strdup(argv[optind]);
  }

  if (opt_show_help) { show_help(); goto cgft_cleanup; }
  if (opt_show_version) { show_version(); goto cgft_cleanup; }

  //if ((opt_create_container + opt_encode + opt_delete + opt_show_band + opt_show_header + opt_show_all) == 0) {
  if ((opt_create_container + opt_encode + opt_delete + opt_show_band + opt_show_header + opt_show_all + opt_run_test + opt_info) == 0) {
    cleanup_ok();
  }

  //if ((opt_create_container + opt_encode + opt_delete + opt_show_band + opt_show_header + opt_show_all) != 1) {
  if ((opt_create_container + opt_encode + opt_delete + opt_show_band + opt_show_header + opt_show_all + opt_run_test + opt_info) != 1) {
    printf("must specify exactly one of show header (-H), show band (-b), encode (-e), delete (-d) or create empty container (-C)\n");
    cleanup_err();
  }

  if (opt_show_header) {
    if (!opt_ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(opt_ifn, "r"))==NULL) { perror(opt_ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", opt_ifn);
      cleanup_fail();
    }

    cgft_print_header(cgf);

  }
  else if (opt_info) {

    if (!opt_ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(opt_ifn, "r"))==NULL) { perror(opt_ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", opt_ifn);
      cleanup_fail();
    }

    cgft_stats(cgf);

  }

  else if (opt_show_band) {
    if (opt_tilepath<0) { printf("must specify tilepath\n"); cleanup_err(); }
    if (!opt_ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(opt_ifn, "r"))==NULL) { perror(opt_ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", opt_ifn);
      cleanup_fail();
    }

    for (idx=0; idx<cgf->Path.size(); idx++) {
      if (cgf->Path[idx].TilePath == (uint64_t)opt_tilepath) { break; }
    }

    if ((uint64_t)idx==cgf->Path.size()) {
      printf("Tile Path %i not found\n", opt_tilepath);
      cleanup_ok();
    }

    cgft_output_band_format(cgf, &(cgf->Path[idx]), stdout);

  }

  else if (opt_encode) {
    if (opt_tilepath<0) { printf("must specify tilepath\n"); cleanup_err(); }
    if (!opt_ifn && opt_ofn)        { opt_ifn=strdup(opt_ofn); }
    else if (opt_ifn && !opt_ofn)   { opt_ofn=strdup(opt_ifn); }
    else if ((!opt_ifn) && (!opt_ofn)) { printf("provide CGF file\n"); cleanup_err(); }

    if ((ifp=fopen(opt_ifn, "r"))==NULL) { perror(opt_ifn); cleanup_err(); }
    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", opt_ifn);
      cleanup_fail();
    }
    if (ifp!=stdin) { fclose(ifp); }
    ifp = NULL;

    if (opt_band_ifn) {
      if ((opt_band_ifp=fopen(opt_band_ifn, "r"))==NULL) { perror(opt_band_ifn); cleanup_err(); }
    } else {
      opt_band_ifp = stdin;
    }

    for (idx=0; idx<cgf->Path.size(); idx++) {
      if (cgf->Path[idx].TilePath == (uint64_t)opt_tilepath) { break; }
    }
    if (idx==cgf->Path.size()) {
      tilepath_t p;
      cgf->Path.push_back(p);
      cgft_tilepath_init(cgf->Path[idx], (uint64_t)opt_tilepath);
      cgf->PathCount++;
    }

    //encode...
    //
    k = cgft_read_band_tilepath(cgf, &(cgf->Path[idx]), opt_band_ifp);

    // Do some rudimentary sanity checks
    //
    k = cgft_sanity(cgf);
    if (k<0) {
      printf("SANITY FAIL: %i\n", k);
      cleanup_fail();
    }

    k = cgft_write_to_file(cgf, opt_ofn);


    printf("got %i\n", k);

  }

  else if (opt_delete) {
    if (opt_tilepath<0) { printf("must specify tilepath\n"); cleanup_err(); }
    if (!opt_ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if (!opt_ofn) { opt_ofn=strdup(opt_ifn); }

    //...
  }

  else if (opt_create_container) {
    //...

    if (!opt_ofn) {
      if (opt_ifn) { opt_ofn=strdup(opt_ifn); }
      else { printf("specify output CGF file\n"); cleanup_err(); }
    }

    if (!opt_tilemap_fn) {
      tilemap_str = DEFAULT_TILEMAP;
    } else {
      if ((read_tilemap_from_file(tilemap_str, opt_tilemap_fn))==NULL) {
        perror(opt_tilemap_fn);
        cleanup_err();
      }
    }

    if (opt_ofn=="-") { ofp=stdout; }
    else if ((ofp = fopen(opt_ofn, "w"))==NULL) { perror(opt_ofn); cleanup_err(); }

    cgft_create_container(ofp, tilemap_str.c_str());

    //cgft_create_container(ofp);
  }

  else if (opt_show_all) {
    if (!opt_ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(opt_ifn, "r"))==NULL) { perror(opt_ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", opt_ifn);
      cleanup_fail();
    }

    printf("... %i\n", (int)cgf->Path.size());

    cgft_print_header(cgf);


    for (idx=0; idx<cgf->Path.size(); idx++) {
      cgft_print_tilepath(cgf, &(cgf->Path[idx]));
    }


  }
  else if (opt_run_test) {

    if (!opt_ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(opt_ifn, "r"))==NULL) { perror(opt_ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", opt_ifn);
      cleanup_fail();
    }

    cgft_run_test(cgf);

    printf("loaded...\n");
  }



cgft_cleanup:
  if (opt_ifn) { free(opt_ifn); }
  if (opt_ofn) { free(opt_ofn); }
  if (opt_band_ifn) { free(opt_band_ifn); }
  if (opt_tilemap_fn) { free(opt_tilemap_fn); }
  if (ifp && (ifp!=stdin)) { fclose(ifp); }
  if (ofp && (ofp!=stdout)) { fclose(ofp); }
  if (opt_band_ifp && (opt_band_ifp!=stdin)) { fclose(opt_band_ifp); }
  exit(ret);
}
