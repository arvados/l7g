#include "cgft.hpp"

#define CGFT_VERSION "0.2.0"

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

typedef struct cgft_opt_type {
  int show_header,
      show_band,
      encode,
      show_help,
      show_version,
      verbose,
      del,
      create_container,
      tilemap,
      show_all,
      ez_print;
  int run_test,
      info;
  char *ifn,
       *ofn,
       *tilemap_fn;
  char *band_ifn;
  FILE *band_ifp;
  int tilepath;

  int cgf_version_opt;
  std::vector< std::string > cgf_version_opt_ele;
  std::string cgf_version_str;
  int update_cgf_version;

  int cglf_version_opt;
  std::vector< std::string > cglf_version_opt_ele;
  std::string cglf_version_str;
  int update_cglf_version;

  int update_header;

} cgft_t;

static cgft_t cgft_opt;

void init_cgft_opt(cgft_t *opt) {
  opt->show_header=0;
  opt->show_band=0;
  opt->encode=0;
  opt->show_help=0;
  opt->show_version=0;
  opt->verbose=0;
  opt->del=0;
  opt->create_container=0;
  opt->tilemap=0;
  opt->show_all=0;
  opt->ez_print=0;
  opt->run_test=0;
  opt->info=0;
  opt->ifn=NULL;
  opt->ofn=NULL;
  opt->tilemap_fn = NULL;
  opt->band_ifn = NULL;
  opt->band_ifp = NULL;
  opt->tilepath = -1;

  opt->cgf_version_opt = 0;
  opt->cgf_version_opt_ele.clear();

  opt->cgf_version_str = CGF_VERSION;
  opt->update_cgf_version = 0;

  opt->cglf_version_opt = 0;
  opt->cglf_version_opt_ele.clear();
  opt->cglf_version_str = CGLF_VERSION;
  opt->update_cglf_version = 0;

  opt->update_header = 0;
}

static struct option long_options[] = {
  {"header",              no_argument,        NULL, 'H'},
  {"create-container",    no_argument,        NULL, 'C'},
  {"info",                no_argument,        NULL, 'I'},
  {"show-all",            no_argument,        NULL, 'A'},
  {"help",                no_argument,        NULL, 'h'},
  {"version",             no_argument,        NULL, 'v'},
  {"verbose",             no_argument,        NULL, 'V'},
  {"ez-print",            no_argument,        NULL, 'Z'},
  {"band",                required_argument,  NULL, 'b'},
  {"encode",              required_argument,  NULL, 'e'},
  {"input",               required_argument,  NULL, 'i'},
  {"output",              required_argument,  NULL, 'o'},
  {"tilemap",             required_argument,  NULL, 't'},
  {"version-opt",         required_argument,  NULL, 'T'},
  {"library-version",     required_argument,  NULL, 'L'},
  {0,0,0,0}
};

void show_help() {
  printf("CGF Tool.  A tool used to inspect and edit Compact Genome Format (CGF) files.\n");
  printf("Version: %s\n", CGFT_VERSION);
  printf("\n");
  //printf("usage:\n");
  //printf("\n");
  printf("usage: cgft [-H] [-b tilepath] [-e tilepath] [-i ifn] [-o ofn] [-h] [-v] [-V] [ifn]\n");
  printf("\n");
  printf("  [-H|--header]               show header\n");
  printf("  [-C|--create-container]     create empty container\n");
  printf("  [-I|--info]                 print basic information about CGF file\n");
  printf("  [-b|--band tilepath]        output band for tilepath\n");
  printf("  [-e|--encode tilepath]      input tilepath band and add it to file, overwriting if it already exists\n");
  printf("  [-i|--input ifn]            input file (CGF)\n");
  printf("  [-o|--output ofn]           output file (CGF)\n");
  printf("  [-A|--show-all]             show all tilepaths\n");
  printf("  [-h|--help]                 show help (this screen)\n");
  printf("  [-v|--version]              show version\n");
  printf("  [-V|--verbose]              set verbose level\n");
  printf("  [-t|--tilemap tfn]          use tilemap file (instead of default)\n");
  printf("  [-Z|--ez-print]             print \"ez\" structure information\n");
  printf("  [-T|--version-opt vopt]     CGF version option.  Must be one of \"default\" or \"noc-inv\"\n");
  printf("  [-L|--library-version lopt] CGF library version option.  Overwrites default value if specified\n");
  printf("  [-U|--update-header]        Update header only\n");
  printf("\n");
}

int main(int argc, char **argv) {
  int i, j, k;
  int opt;
  int ret=0;
  int idx;
  char buf[1024];

  std::string tilemap_str;
  FILE *ifp=NULL, *ofp=NULL;
  cgf_t *cgf=NULL;
  int option_index=0;
  int def_or_nocinv=0;

  init_cgft_opt(&cgft_opt);

  while ((opt = getopt_long(argc, argv, "Hb:e:i:o:Ct:T:L:U:hvVAZRI", long_options, &option_index))!=-1) switch (opt) {
    case 0:
      fprintf(stderr, "sanity error, invalid option to parse, exiting\n");
      exit(-1);
      break;
    case 'H': cgft_opt.show_header=1; break;
    case 'C': cgft_opt.create_container=1; break;
    case 'I': cgft_opt.info=1; break;
    case 't': cgft_opt.tilemap=1; cgft_opt.tilemap_fn=strdup(optarg); break;
    case 'b': cgft_opt.show_band=1; cgft_opt.tilepath=atoi(optarg); break;
    case 'e': cgft_opt.encode=1; cgft_opt.tilepath=atoi(optarg); break;
    case 'd': cgft_opt.del=1; cgft_opt.tilepath=atoi(optarg); break;
    case 'i': cgft_opt.ifn=strdup(optarg); break;
    case 'o': cgft_opt.ofn=strdup(optarg); break;
    case 'h': cgft_opt.show_help=1; break;
    case 'A': cgft_opt.show_all=1; break;
    case 'v': cgft_opt.show_version=1; break;
    case 'V': cgft_opt.verbose=1; break;
    case 'Z': cgft_opt.ez_print=1; break;
    case 'R': cgft_opt.run_test=1; break;
    case 'T': cgft_opt.cgf_version_opt=1;
              cgft_opt.cgf_version_opt_ele.push_back(optarg);
              cgft_opt.update_header = 1;
              break;
    case 'L': cgft_opt.cglf_version_opt=1;
              cgft_opt.cglf_version_opt_ele.push_back(optarg);
              cgft_opt.update_header=1;
              break;
    case 'U': cgft_opt.update_header=1; break;
    default: printf("unknown option"); show_help(); cleanup_ok(); break;
  }

  if (argc>optind) {
    if ((argc-optind)>1) { printf("Extra options specified\n"); cleanup_err(); }
    if (cgft_opt.ifn) { printf("Input CGF already specified.\n"); cleanup_err(); }
    cgft_opt.ifn = strdup(argv[optind]);
  }

  if (cgft_opt.show_help) { show_help(); goto cgft_cleanup; }
  if (cgft_opt.show_version) { show_version(); goto cgft_cleanup; }

  // We must have a command.  If not, exit.
  // The 'no-surprise' rule is to show help and exit gracefully when
  // no commands are specified.
  //
  if ((cgft_opt.create_container +
       cgft_opt.encode +
       cgft_opt.del +
       cgft_opt.show_band +
       cgft_opt.show_header +
       cgft_opt.show_all +
       cgft_opt.run_test +
       cgft_opt.info +
       cgft_opt.update_header) == 0) {
    cleanup_ok();
  }

  // Don't allow more than one command
  //
  if ((cgft_opt.create_container +
       cgft_opt.encode +
       cgft_opt.del +
       cgft_opt.show_band +
       cgft_opt.show_header +
       cgft_opt.show_all +
       cgft_opt.run_test +
       cgft_opt.info) > 0) {
    cgft_opt.update_header = 0;
  }

  // Don't allow more than one command
  //
  if ((cgft_opt.create_container +
       cgft_opt.encode +
       cgft_opt.del +
       cgft_opt.show_band +
       cgft_opt.show_header +
       cgft_opt.show_all +
       cgft_opt.run_test +
       cgft_opt.info +
       cgft_opt.update_header) != 1) {
    printf("must specify exactly one of show header (-H), show band (-b), encode (-e), delete (-d), create empty container (-C) or update header (-U)\n");
    cleanup_err();
  }

  // Create version string to save
  //
  for (i=0; i<cgft_opt.cgf_version_opt_ele.size(); i++) {
    cgft_opt.update_cgf_version = 1;
    if ( strncmp(cgft_opt.cgf_version_opt_ele[i].c_str(), "default", strlen("default"))==0 ) {
      def_or_nocinv=0;
    }
    else if ( strncmp(cgft_opt.cgf_version_opt_ele[i].c_str(), "noc-inv", strlen("noc-inv"))==0 ) {
      def_or_nocinv=1;
    }
    else {
      printf("invalid CGFVersion option.  Must be one of 'default' or 'noc-inv'.\n");
      cleanup_err();
    }
  }
  if (def_or_nocinv==1) {
    cgft_opt.cgf_version_str += ",";
    cgft_opt.cgf_version_str += "noc-inv";
  }

  // Create CGF library version string to save.
  // Overwrite default value if we have at least one option.
  //
  for (i=0; i<cgft_opt.cglf_version_opt_ele.size(); i++) {
    cgft_opt.update_cglf_version = 1;
    if (i==0) { cgft_opt.cglf_version_str.clear(); }
    else { cgft_opt.cglf_version_str += ","; }
    cgft_opt.cglf_version_str += cgft_opt.cglf_version_opt_ele[i];
  }

  //-------------------------
  //
  // Process command
  //
  //-------------------------

  // Show header
  //

  if (cgft_opt.show_header) {
    if (!cgft_opt.ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }

    cgft_print_header(cgf);

  }

  // CGF info (some simple stats)
  //

  else if (cgft_opt.info) {

    if (!cgft_opt.ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }

    cgft_stats(cgf);

  }

  // Show 'band' output for a tile path
  //

  else if (cgft_opt.show_band) {
    if (cgft_opt.tilepath<0) { printf("must specify tilepath\n"); cleanup_err(); }
    if (!cgft_opt.ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }

    for (idx=0; idx<cgf->Path.size(); idx++) {
      if (cgf->Path[idx].TilePath == (uint64_t)cgft_opt.tilepath) { break; }
    }

    if ((uint64_t)idx==cgf->Path.size()) {
      printf("Tile Path %i not found\n", cgft_opt.tilepath);
      cleanup_ok();
    }

    cgft_output_band_format(cgf, &(cgf->Path[idx]), stdout);

  }

  else if (cgft_opt.update_header) {

    if (!cgft_opt.ifn && cgft_opt.ofn)        { cgft_opt.ifn=strdup(cgft_opt.ofn); }
    else if (cgft_opt.ifn && !cgft_opt.ofn)   { cgft_opt.ofn=strdup(cgft_opt.ifn); }
    else if ((!cgft_opt.ifn) && (!cgft_opt.ofn)) { printf("provide CGF file\n"); cleanup_err(); }

    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }
    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }
    if (ifp!=stdin) { fclose(ifp); }
    ifp = NULL;

    // Update header information if necessary
    //
    if (cgft_opt.update_header) {
      cgf->CGFVersion = cgft_opt.cgf_version_str;
      cgf->LibraryVersion= cgft_opt.cglf_version_str;
    }

    // Do some rudimentary sanity checks
    //
    k = cgft_sanity(cgf);
    if (k<0) {
      fprintf(stderr, "SANITY FAIL: %i\n", k);
      cleanup_fail();
    }

    k = cgft_write_to_file(cgf, cgft_opt.ofn);

  }

  // Encode a tile path from 'band' representation input
  //

  else if (cgft_opt.encode) {
    if (cgft_opt.tilepath<0) { printf("must specify tilepath\n"); cleanup_err(); }
    if (!cgft_opt.ifn && cgft_opt.ofn)        { cgft_opt.ifn=strdup(cgft_opt.ofn); }
    else if (cgft_opt.ifn && !cgft_opt.ofn)   { cgft_opt.ofn=strdup(cgft_opt.ifn); }
    else if ((!cgft_opt.ifn) && (!cgft_opt.ofn)) { printf("provide CGF file\n"); cleanup_err(); }

    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }
    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }
    if (ifp!=stdin) { fclose(ifp); }
    ifp = NULL;

    // Update header information if necessary
    //
    if (cgft_opt.update_header) {
      cgf->CGFVersion = cgft_opt.cgf_version_str;
      cgf->LibraryVersion= cgft_opt.cglf_version_str;
    }

    if (cgft_opt.band_ifn) {
      if ((cgft_opt.band_ifp=fopen(cgft_opt.band_ifn, "r"))==NULL) { perror(cgft_opt.band_ifn); cleanup_err(); }
    } else {
      cgft_opt.band_ifp = stdin;
    }

    for (idx=0; idx<cgf->Path.size(); idx++) {
      if (cgf->Path[idx].TilePath == (uint64_t)cgft_opt.tilepath) { break; }
    }
    if (idx==cgf->Path.size()) {
      tilepath_t p;
      cgf->Path.push_back(p);
      cgft_tilepath_init(cgf->Path[idx], (uint64_t)cgft_opt.tilepath);
      cgf->PathCount++;
    }

    //encode...
    //
    k = cgft_read_band_tilepath(cgf, &(cgf->Path[idx]), cgft_opt.band_ifp);

    // Do some rudimentary sanity checks
    //
    k = cgft_sanity(cgf);
    if (k<0) {
      fprintf(stderr, "SANITY FAIL: %i\n", k);
      cleanup_fail();
    }

    k = cgft_write_to_file(cgf, cgft_opt.ofn);

  }

  // Delete a tile path.
  // There probably should be a difference between 'delete' and 'remove'.
  // 'delete' should clear out all relevant data but still leave a place holder (?)
  // 'remove' should remove the path entirely.  It's a little unclear what this means
  //   if the tile path is in the middle of a contiguous region.
  //

  else if (cgft_opt.del) {
    if (cgft_opt.tilepath<0) { printf("must specify tilepath\n"); cleanup_err(); }
    if (!cgft_opt.ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if (!cgft_opt.ofn) { cgft_opt.ofn=strdup(cgft_opt.ifn); }

    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }
    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }
    if (ifp!=stdin) { fclose(ifp); }
    ifp = NULL;

    // Update header information if necessary
    //
    if (cgft_opt.update_header) {
      cgf->CGFVersion = cgft_opt.cgf_version_str;
      cgf->LibraryVersion = cgft_opt.cglf_version_str;
    }

    for (idx=0; idx<cgf->Path.size(); idx++) {
      if (cgf->Path[idx].TilePath == (uint64_t)cgft_opt.tilepath) { break; }
    }

    if (idx==cgf->Path.size()) {
      fprintf(stderr, "Path not found (%i %x)\n", cgft_opt.tilepath, cgft_opt.tilepath);
      cleanup_fail();
    }
    else {
      tilepath_t null_tilepath;
      cgf->Path[idx] = null_tilepath;
      cgft_tilepath_init(cgf->Path[idx], (uint64_t)cgft_opt.tilepath);
    }

    // Do some rudimentary sanity checks
    //
    k = cgft_sanity(cgf);
    if (k<0) {
      fprintf(stderr, "SANITY FAIL: %i\n", k);
      cleanup_fail();
    }

    k = cgft_write_to_file(cgf, cgft_opt.ofn);


    //...
  }

  // Create an empty CGF container.
  //

  else if (cgft_opt.create_container) {

    if (!cgft_opt.ofn) {
      if (cgft_opt.ifn) { cgft_opt.ofn=strdup(cgft_opt.ifn); }
      else { printf("specify output CGF file\n"); cleanup_err(); }
    }

    if (!cgft_opt.tilemap_fn) {
      tilemap_str = DEFAULT_TILEMAP;
    } else {
      if ((read_tilemap_from_file(tilemap_str, cgft_opt.tilemap_fn))==NULL) {
        perror(cgft_opt.tilemap_fn);
        cleanup_err();
      }
    }

    if (cgft_opt.ofn=="-") { ofp=stdout; }
    else if ((ofp = fopen(cgft_opt.ofn, "w"))==NULL) { perror(cgft_opt.ofn); cleanup_err(); }

    cgft_create_container(
        ofp,
        cgft_opt.cgf_version_str.c_str(),
        cgft_opt.cglf_version_str.c_str(),
        tilemap_str.c_str());
  }

  // Print out a verbose representation of the contenst of the CGF file
  //

  else if (cgft_opt.show_all) {
    if (!cgft_opt.ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }

    printf("... %i\n", (int)cgf->Path.size());

    cgft_print_header(cgf);

    for (idx=0; idx<cgf->Path.size(); idx++) {
      cgft_print_tilepath(cgf, &(cgf->Path[idx]));
    }

  }

  // Run some tests on the CGF file
  //

  else if (cgft_opt.run_test) {

    if (!cgft_opt.ifn) { printf("provide input CGF file\n"); cleanup_err(); }
    if ((ifp=fopen(cgft_opt.ifn, "r"))==NULL) { perror(cgft_opt.ifn); cleanup_err(); }

    cgf = cgft_read(ifp);
    if (!cgf) {
      printf("CGF read error.  Is %s a valid CGFv3 file?\n", cgft_opt.ifn);
      cleanup_fail();
    }

    cgft_run_test(cgf);

    printf("loaded...\n");
  }



cgft_cleanup:
  if (cgft_opt.ifn) { free(cgft_opt.ifn); }
  if (cgft_opt.ofn) { free(cgft_opt.ofn); }
  if (cgft_opt.band_ifn) { free(cgft_opt.band_ifn); }
  if (cgft_opt.tilemap_fn) { free(cgft_opt.tilemap_fn); }
  if (ifp && (ifp!=stdin)) { fclose(ifp); }
  if (ofp && (ofp!=stdout)) { fclose(ofp); }
  if (cgft_opt.band_ifp && (cgft_opt.band_ifp!=stdin)) { fclose(cgft_opt.band_ifp); }
  exit(ret);
}
