#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>
#include <errno.h>

#include <htslib/hts.h>
#include <htslib/vcf.h>
#include <htslib/faidx.h>
#include <htslib/kstring.h>
#include <htslib/synced_bcf_reader.h>

#include <vector>
#include <string>

#define VCFBED2HOMREF_VERSION "0.1.3"

int verbose_flag = 0;
int show_header_flag = 1;

std::string NON_REF_STR = "<NON_REF>";

void show_version(FILE *ofp) {
  fprintf(ofp, "vcfbed2homref version: %s\n", VCFBED2HOMREF_VERSION);
}

void show_help(FILE *ofp) {
  show_version(ofp);
  fprintf(ofp, "usage:\n\n  vcfbed2homref [-h] [-v] [-V] [-N non-ref] [-b bedfile] [-r ref-fasta] [vcf_file] [out_vcf_file]\n");
  fprintf(ofp, "\n");
  fprintf(ofp, "    [vcf_file]      VCF file (defaults stdin)\n");
  fprintf(ofp, "    [out_vcf_file]  output VCF file (defaults stdout)\n");
  fprintf(ofp, "    [-b bedfile]    bed file of homozygous ref sequences\n");
  fprintf(ofp, "    [-r ref-fasta]  reference FASTA file (indexed)\n");
  fprintf(ofp, "    [-s]            supress header (default output header)\n");
  fprintf(ofp, "    [-v]            verbose\n");
  fprintf(ofp, "    [-N non-ref]    \"non ref\" string to use (defaults to '%s')\n", NON_REF_STR.c_str());
  fprintf(ofp, "    [-V]            print version\n");
  fprintf(ofp, "    [-h]            Help (this screen)\n");
  fprintf(ofp, "\n");
}

static struct option long_options[] = {
  {"help", no_argument, NULL, 'h'},
  {"verbose", no_argument, NULL, 'v'},
  {"version", no_argument, NULL, 'V'},
  {"no-header", no_argument, NULL, 's'},
  {"bed", required_argument, NULL, 'b'},
  {"ref", required_argument, NULL, 'r'},
  {"input", required_argument, NULL, 'i'},
  {"output", required_argument, NULL, 'o'},
  {"non-ref-region", required_argument, NULL, 'N'},
  {0,0,0,0}
};

typedef struct bed_info_type {
  std::vector< std::string > block_chrom;
  std::vector< int > index;

  std::vector< std::string > chrom;
  std::vector< int > start0ref;
  std::vector< int > end0noninc;
} bed_info_t;


int parse_bed_line(std::string &line, std::string &chrom, int *s0, int *e0) {
  int i, idx=0, n;
  long long int v;

  n = (int)line.size();
  chrom.clear();
  while ((idx<n) && (line[idx]!='\t')) {
    chrom += line[idx];
    idx++;
  }
  idx++;
  if (idx>=n) { return -1; }

  v = strtoll( (&(line[0])) + idx, NULL, 10);
  while ((idx<n) && (line[idx]!='\t')) {
    if ((line[idx]<'0') || (line[idx]>'9')) { return -2; }
    idx++;
  }
  idx++;
  if (idx>=n) { return -2; }

  *s0 = (int)v;

  v = strtoll( (&(line[0])) + idx, NULL, 10);
  while ((idx<n) && (line[idx]!='\t') && (line[idx]!='\n')) {
    if ((line[idx]<'0') || (line[idx]>'9')) { return -2; }
    idx++;
  }

  *e0 = (int)v;

  return 0;

}

int read_bed(FILE *fp, bed_info_t &bed_info) {
  std::string buf, chrom;
  int s0, e0;
  int ch, ret=0;
  int line_no = 0;

  std::string block_chrom;
  int idx=0;

  while (!feof(fp)) {
    ch = fgetc(fp);
    if ((ch==EOF) || (ch=='\n')) {
      if (buf.size()==0) { continue; }
      ret = parse_bed_line(buf, chrom, &s0, &e0);
      if (ret<0) { return ret; }
      buf.clear();

      bed_info.chrom.push_back(chrom);
      bed_info.start0ref.push_back(s0);
      bed_info.end0noninc.push_back(e0);

      if (idx>0) {
        if (block_chrom != chrom) {
          bed_info.block_chrom.push_back(block_chrom);
          bed_info.index.push_back(idx);
        }
      }
      idx = bed_info.chrom.size();
      block_chrom = chrom;

      line_no++;
      continue;
    }
    buf += (char)ch;
  }
  bed_info.block_chrom.push_back(block_chrom);
  bed_info.index.push_back(idx);

  return 0;
}

void print_vcf_homref(FILE *fp, const char *chrom, int start0ref, int n, const char *ref_anchor, int n_samples) {
  int i, phased=0;
  fprintf(fp, "%s\t%i\t%s\t%s\t%s\t%s\t%s\tEND=%i\t%s",
      chrom,
      start0ref+1,
      ".",
      ref_anchor,
      NON_REF_STR.c_str(),
      ".",
      "PASS",
      start0ref+n,
      "GT");
  for (i=0; i<n_samples; i++) { fprintf(fp, "\t0%c0", phased ? '|' : '/' ); }
  fprintf(fp, "\n");
}

int faidx_seq(faidx_t *fai, const char *chrom, int s0, int n, std::string &res) {
  int len=0;
  char *s, reg[1024];
  sprintf(reg, "%s:%i-%i", chrom, s0+1, s0+n);
  s = fai_fetch(fai, reg, &len);
  if (!s) { return -1; 
  fprintf(stderr, "No fai file found...\n"); }
  if (len<0) { free(s); return -2; }

  res = s;

  free(s);
  return 0;
}


int process_vcf2homref(std::string &ofn, std::string &vcf_fn, bed_info_t &bed_info, std::string &ref_fn) {
  int i, r, ret=0;
  faidx_t *fai;
  FILE *ofp;

  std::string chrom_name;

  //htsFile *ofp;
  kstring_t ks={0,0,0};

  htsFile *bcf_fp = NULL;
  bcf_hdr_t *bcf_hdr = NULL;
  bcf1_t *bcf_line = NULL;
  bcf_srs_t *bcf_stream = NULL;

  int homref_beg = 0, bed_idx = 0, bed_end_idx = 0;
  int vcf_beg_pos = 0;

  int count=0,ncount=10;
  int nref=0, n_sample=0;
  std::string ref_anchor;

  if (ofn=="") { ofp = stdout; }
  else {
    ofp = fopen(ofn.c_str(), "w");
    if (!ofp) { return -1; 
    fprintf(stderr, "Can't open input file (maybe bed?)\n");
    }
  }


  homref_beg = 0;
  bed_idx = 0;
  chrom_name = "unk";

  fai = fai_load(ref_fn.c_str());
  if (!fai) { return -3; }

  bcf_fp = bcf_open(vcf_fn.c_str(), "r");
  if (bcf_fp== NULL) { return -1; fprintf(stderr, "Can't open binary vcf...");}

  bcf_hdr = bcf_hdr_read(bcf_fp);
  if (bcf_hdr == NULL) { return -2; }

  // initalize vcf/bcf streams
  //
  bcf_stream = bcf_sr_init();
  bcf_sr_add_reader(bcf_stream, vcf_fn.c_str());

  bed_idx = 0;
  bed_end_idx = 0;

  if (show_header_flag) {
    r = bcf_hdr_format(bcf_hdr, 0, &ks);
    if (r<0) {
      ret = -1;
      fprintf(stderr, "Can't read bcf header...\n");
      goto process_vcfbed2homref_cleanup;
    }
    fwrite(ks.s, sizeof(char), ks.l, ofp);
    ks.l=0;
  }

  while (bcf_sr_next_line(bcf_stream)) {

    // We're only processing 1 file, '0' here is the 0th (and only)
    // file.
    //
    bcf_line = bcf_sr_get_line(bcf_stream, 0);

    chrom_name = bcf_hdr_id2name(bcf_hdr, bcf_line->rid);
    vcf_beg_pos = bcf_line->pos;
    nref = (int)bcf_line->rlen;
    n_sample = (int)bcf_line->n_sample;

    // If we've changed from one chromosome block
    // to another in the VCF file, emit all
    // remaining homref lines implied by the BED file.
    //
    if ( (bed_idx < bed_end_idx) &&
         (chrom_name != bed_info.chrom[bed_idx]) ) {

      while (bed_idx < bed_end_idx) {

        r = faidx_seq(fai, bed_info.chrom[bed_idx].c_str(), homref_beg, 1, ref_anchor);
        if (r<0) {
          ret = -14;
          goto process_vcfbed2homref_cleanup;
        }

        print_vcf_homref(ofp,
            bed_info.chrom[bed_idx].c_str(),
            homref_beg,
            bed_info.end0noninc[bed_idx] - homref_beg,
            ref_anchor.c_str(),
            n_sample);
        bed_idx++;
        if (bed_idx < bed_end_idx) {
          homref_beg = bed_info.start0ref[bed_idx];
        }
      }

    }

    // If we've reached the end of our chromosome block
    // in the VCF (or we've just started), search for
    // the bed range and initialize state.
    //
    if (bed_idx >= bed_end_idx) {

      for (i=0; i<bed_info.block_chrom.size(); i++) {
        if (bed_info.block_chrom[i]==chrom_name) {
          bed_idx = ((i==0) ? 0 : bed_info.index[i-1]);
          bed_end_idx = bed_info.index[i];
          homref_beg = bed_info.start0ref[bed_idx];
          break;
        }
      }

      // we didn't find the chromosome, return with error
      //
      if (i==bed_info.block_chrom.size()) {
        ret = -1;
        fprintf(stderr, "Can't find the chromosome...\n");
        goto process_vcfbed2homref_cleanup;
      }

    }

    // advance the bed information to currently read beginning of vcf file
    //
    if (homref_beg < vcf_beg_pos) {

      // there could be multiple BED entries, so advance until the bed
      // end range falls past the beginning of the vcf range
      //
      while ( (bed_idx < bed_end_idx) &&
              (chrom_name == bed_info.chrom[bed_idx]) &&
              (bed_info.end0noninc[bed_idx] <= vcf_beg_pos) ) {

        r = faidx_seq(fai, bed_info.chrom[bed_idx].c_str(), homref_beg, 1, ref_anchor);
        if (r<0) {
          ret = -12;
          goto process_vcfbed2homref_cleanup;
        }

        print_vcf_homref(ofp,
            bed_info.chrom[bed_idx].c_str(),
            homref_beg,
            bed_info.end0noninc[bed_idx] - homref_beg,
            ref_anchor.c_str(),
            n_sample);


        bed_idx++;
        if ( (bed_idx < bed_end_idx) &&
             (chrom_name == bed_info.chrom[bed_idx]) ) {
          homref_beg = bed_info.start0ref[bed_idx];
        }
      }

      // at this point, the bed range end should fall past the vcf beginning.
      // If the current `homref_beg` range is less than the reported
      // VCF beginning, the current BED range straddles the vcf, so
      // print the beginning hom. ref. window.
      //
      if ( (bed_idx < bed_end_idx) &&
           (chrom_name == bed_info.chrom[bed_idx]) &&
           (homref_beg < vcf_beg_pos) ) {

        r = faidx_seq(fai, bed_info.chrom[bed_idx].c_str(), homref_beg, 1, ref_anchor);
        if (r<0) {
          ret = -11;
          goto process_vcfbed2homref_cleanup;

        }

        print_vcf_homref(ofp,
            bed_info.chrom[bed_idx].c_str(),
            homref_beg,
            vcf_beg_pos - homref_beg,
            ref_anchor.c_str(),
            n_sample);

      }

    }

    vcf_format(bcf_hdr, bcf_line, &ks);
    fwrite(ks.s, sizeof(char), ks.l, ofp);
    ks.l=0;

    // if homref_beg is less than the end of the current
    // reference length reported by the vcf line,
    // update it's position. The homref block should
    // have already been printed at this point in this
    // case.
    //
    if (homref_beg < (vcf_beg_pos + nref)) {
      homref_beg = vcf_beg_pos + nref;
    }

    // Now that we've advanced the `homref_beg` variable,
    // we want to update the `bed_info` index position as well,
    // making sure the end of the `bed_info` doesn't fall
    // before where we are in our homref position.
    //
    while ( (bed_idx < bed_end_idx) &&
            (chrom_name == bed_info.chrom[bed_idx]) &&
            (bed_info.end0noninc[bed_idx] <= homref_beg) ) {
      bed_idx++;
    }

    // Once we've updated the `bed_info` index, update
    // the `homref_beg` variable if it falls before
    // the current `bed_info` window.
    //
    if ( (bed_idx < bed_end_idx) &&
         (chrom_name == bed_info.chrom[bed_idx]) &&
         (homref_beg < bed_info.start0ref[bed_idx]) ) {
      homref_beg = bed_info.start0ref[bed_idx];
    }

  }

  // We've reached the end of the file but we still have
  // some homozygous reference ntries to take care of,
  // so take care of them here.
  //
  if (bed_idx < bed_end_idx) {

    while (bed_idx < bed_end_idx) {

      r = faidx_seq(fai, bed_info.chrom[bed_idx].c_str(), homref_beg, 1, ref_anchor);
      if (r<0) {
        ret = -10;
        goto process_vcfbed2homref_cleanup;
      }

      print_vcf_homref(ofp,
          bed_info.chrom[bed_idx].c_str(),
          homref_beg,
          bed_info.end0noninc[bed_idx] - homref_beg,
          ref_anchor.c_str(),
          n_sample);
      bed_idx++;
      if (bed_idx < bed_end_idx) {
        homref_beg = bed_info.start0ref[bed_idx];
      }
    }

  }

process_vcfbed2homref_cleanup:

  fai_destroy(fai);

  free(ks.s);

  bcf_sr_destroy(bcf_stream);
  bcf_hdr_destroy(bcf_hdr);
  bcf_close(bcf_fp);

  if (ofp!=stdout) { fclose(ofp); }

  return ret;
}

int main(int argc, char **argv) {
  int i, j, k, n, ret=0, ch;
  int opt, option_index;
  bed_info_t bed_info;

  std::string bed_fn = ""; //  "./small.bed";
  std::string ref_fn = ""; // "/data-sdd/data/ref/human_g1k_v37.fa.gz";

  std::string vcf_fn = "", out_vcf_fn = "";

  FILE *bed_fp = stdin;

  while ((opt=getopt_long(argc, argv, "vVhb:r:i:o:sN:", long_options, &option_index))!=-1) switch(opt) {
    case 0:
      fprintf(stderr, "invalid option, exiting\n");
      exit(-1);
      break;
    case 'h':
      show_help(stdout);
      exit(0);
      break;
    case 'V':
      show_version(stdout);
      exit(0);
      break;
    case 'v':
      verbose_flag=1;
      break;

    case 'b':
      bed_fn = optarg;
      break;
    case 'r':
      ref_fn = optarg;
      break;
    case 's':
      show_header_flag = 0;
      break;

    case 'N':
      NON_REF_STR = optarg;
      break;

    case 'i':
      vcf_fn = optarg;
      break;
    case 'o':
      out_vcf_fn = optarg;
      break;
    default:
      show_help(stderr);
      exit(-1);
      break;
  }

  if (optind < argc) {
    if ( optind == (argc-2) ) {
      vcf_fn = argv[optind];
      out_vcf_fn = argv[optind+1];
    }
    else if (optind == (argc-1)) {
      vcf_fn = argv[optind];
    }
    else {
      fprintf(stderr, "invalid number of options specified\n");
      show_help(stderr);
      exit(-1);
    }
  }

  if ((bed_fn == "") && (vcf_fn == "")) {
    show_help(stderr);
    exit(-1);
  }

  if (vcf_fn == "") { vcf_fn = "/dev/stdin"; }
  else if (bed_fn == "") { bed_fn = "/dev/stdin"; }
  if (out_vcf_fn == "") { out_vcf_fn = "/dev/stdout"; }

  if ((bed_fn == "") || (vcf_fn == "") || (ref_fn == "")) {
    fprintf(stderr, "must provide BED file, VCF file and reference FASTA file\n");
    show_help(stderr);
    exit(-1);
  }

  /*
  printf(">>> bed %s, ref %s, vcf %s, out %s\n",
      bed_fn.c_str(),
      ref_fn.c_str(),
      vcf_fn.c_str(),
      out_vcf_fn.c_str());
      */


  bed_fp = fopen(bed_fn.c_str(), "r");
  if (bed_fp==NULL) {
    perror(bed_fn.c_str());
    exit(-1);
  }
  read_bed(bed_fp, bed_info);
  fclose(bed_fp);

  ret = process_vcf2homref(out_vcf_fn, vcf_fn, bed_info, ref_fn);
  if (ret<0) {
    fprintf(stderr, "ERROR: got return code %i\n", ret);
    exit(-1);
  }
  return 0;
}

