#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>

#include "asm_ukk.h"

#include <string>
#include <vector>

#define WHICH_REF_VERSION "0.1.0"

typedef struct opt_type {
  int debug;
  int show_all_score;
  int verbose;
  int version;
} opt_t;

static struct option long_options[] = {
  {"print-index",         no_argument,        NULL, 'N'},
  {"print-file",          no_argument,        NULL, 'S'},
  {"1ref",                no_argument,        NULL, '1'},
  {"base-concordance",    no_argument,        NULL, 'C'},
  {"case-insensitive",    no_argument,        NULL, 'M'},
  {"help",                no_argument,        NULL, 'h'},
  {"version",             no_argument,        NULL, 'v'},
  {"verbose",             no_argument,        NULL, 'V'},
  {0,0,0,0}
};


void print_usage() {
  printf("\n");
  printf("usage:\n");
  printf("\n");
  printf("    which-ref [-h] [-v] [-V] [query-seq0] [query-seq1] [query-seq2] ... [query-seqN] [ref-seq]\n");
  printf("\n");
  printf("  [-N]        print index (0 reference) of found query seq only\n");
  printf("  [-1]        1 reference (default 0, for base concordance only)\n");
  printf("  [-C]        do base concordance instead of raw sequence concordance\n");
  printf("  [-M]        case insensitive\n");
  printf("  [-S]        print provided filename of found sequence only (overrides index print)\n");
  printf("  [-h]        Print help (this screen)\n");
  printf("  [-v]        Print version\n");
  printf("  [-V]        Verbose flag\n");
  printf("\n");
}

void print_version() {
  printf("version: %s\n", WHICH_REF_VERSION);
}

void init_opt(opt_t *opt) {
  opt->debug = 0;
  opt->show_all_score = 0;
  opt->verbose = 0;
}

int split_tok_ch(std::vector< std::string > &tok_v, std::string &line, const char *fields) {
  int i, j, k, n, n_fields;
  std::string buf;
  char default_fields[] = "\t";
  int subsume_whitespace = 1;

  n = (int)line.size();
  tok_v.clear();

  if (!fields) {
    fields = default_fields;
  }
  n_fields = strlen(fields);

  for (i=0; i<n; i++) {
    for (j=0; j<n_fields; j++) {
      if (line[i]==fields[j]) {
        if (buf.size()>0) { tok_v.push_back(buf); }
        buf.clear();
        break;
      }
    }

    if (j<n_fields) {

      if (subsume_whitespace) {
        while ((i+1)<n) {
          for (j=0; j<n_fields; j++) {
            if (line[i+1] == fields[j]) { break; }
          }
          if (j==n_fields) { break; }
          i++;
        }

      }

      continue;
    }
    buf += line[i];
  }
  if (buf.size()>0) { tok_v.push_back(buf); }

  return 0;
}


typedef struct score_base_type {
  int pos;
  int base;
} score_base_t;

typedef struct score_type {
  int match;
  int mismatch;
  int total;
} score_t;

int _tol(int ch) {
  if ((ch>='A') && (ch<='Z')) {
    return ch - 'A' + 'a';
  }
  return ch;
}

int do_base_concordance(std::string &ifn, std::vector< std::string > &ref_ifns,
                     double *max_score, int *max_idx,
                     int coordinate_0ref, int case_insensitive) {
  int i, j, k;
  int ch;
  FILE *ifp, *ref_fp;
  std::vector< score_t > scores;
  std::vector< std::string > tok;
  std::vector< score_base_t > input_v;
  std::string buf;
  int line_no=0, query_idx, ref_idx;
  score_t s3;

  int loc_debug = 0;

  score_base_t sbt;

  ifp = fopen(ifn.c_str(), "r");
  if (!ifp) { return -1; }

  // read in query pos-base into memory
  //
  buf.clear();
  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if ((ch=='\n') || (ch==EOF)) {

      split_tok_ch(tok, buf, "\t ");
      if (tok.size()!=2) { continue; }

      sbt.pos = atoi(tok[0].c_str());
      sbt.base = ( case_insensitive ? _tol((int)(tok[1][0])) : (int)(tok[1][0]) );
      input_v.push_back(sbt);

      line_no++;
      buf.clear();
      continue;
    }
    buf += (char)ch;
  }
  fclose(ifp);

  //DEBUG
  if (loc_debug) {
    for (i=0; i<input_v.size(); i++) {
      printf("[%i] %i %c\n",
          i, input_v[i].pos, (char)input_v[i].base);
    }
  }

  // Stuff in score structures into our score array
  //
  s3.match=0;
  s3.mismatch=0;
  s3.total=0;
  for (i=0; i<ref_ifns.size(); i++) {
    scores.push_back(s3);
  }

  // Go through each of the references, reading one line
  // at a time (for each) and comparing them to our
  // in-memory query pos-bases.
  //
  for (ref_idx=0; ref_idx<ref_ifns.size(); ref_idx++) {
    ref_fp = fopen(ref_ifns[ref_idx].c_str(), "r");
    if (!ref_fp) { return -1; }

    query_idx = 0;

    line_no=0;
    buf.clear();
    while (!feof(ref_fp)) {

      if (query_idx >= input_v.size()) { break; }

      ch=fgetc(ref_fp);
      if ((ch=='\n') || (ch==EOF)) {

        split_tok_ch(tok, buf, "\t ");
        if (tok.size()!=2) { continue; }

        sbt.pos = atoi(tok[0].c_str());
        sbt.base = ( case_insensitive ? _tol((int)(tok[1][0])) : (int)(tok[1][0]) );

        while ((query_idx < input_v.size()) &&
               (input_v[query_idx].pos < sbt.pos)) {
          query_idx++;
        }
        if (query_idx >= input_v.size()) { break; }

        if (input_v[query_idx].pos == sbt.pos) {
          scores[ref_idx].total++;
          if (input_v[query_idx].base == sbt.base) {

            if (loc_debug) {
              printf("# MATCH @ query[%i] %i %c, ref[%i] %i %c\n",
                  query_idx, input_v[query_idx].pos, (char)(input_v[query_idx].base),
                  ref_idx, sbt.pos, (char)sbt.base);
            }

            scores[ref_idx].match++;
          } else {

            if (loc_debug) {
              printf("# mismatch @ query[%i] %i %c, ref[%i] %i %c\n",
                  query_idx, input_v[query_idx].pos, (char)(input_v[query_idx].base),
                  ref_idx, sbt.pos, (char)sbt.base);
            }

            scores[ref_idx].mismatch++;
          }
        }


        line_no++;
        buf.clear();

        continue;
      }
      buf += (char)ch;
    }

    fclose(ref_fp);
  }

  int processed_first_min_score = 0;
  double d, d_min_score=0.0, d_max_score=0;
  for (i=0; i<scores.size(); i++) {
    if (scores[i].total > 0) {
      d = (double)scores[i].match / (double)scores[i].total;
      if ((!processed_first_min_score) ||
          (d>d_max_score)) {
        d_max_score = d;
        *max_idx = i;
      }
      processed_first_min_score = 1;
    }
  }
  *max_score = d_max_score;

  if (loc_debug) {
    for (i=0; i<scores.size(); i++) {
      printf("score[%i] m:%i, mm:%i, t:%i (%f) %c\n",
          i,
          scores[i].match,
          scores[i].mismatch,
          scores[i].total,
          ((scores[i].total > 0) ? (float)((double)scores[i].match / (double)scores[i].total) : 0.0),
          ( (i==(*max_idx)) ? '*' : ' ') );
    }
  }


  return 0;

}

int main(int argc, char **argv) {
  FILE *fp;
  int i, k, r;
  std::vector< std::string > ref_fns, ref_seq;
  std::string ifn, seq, tseq;
  int loc_debug = 0;

  int print_filename_flag = 0,
      print_index_flag = 0;

  int base_concordance=0,
      coordinate_0ref=1,
      case_insensitive=0;

  int ch, option_index;
  opt_t opt;

  int min_score=-1;
  double max_score=-1;
  int max_idx, min_idx;

  std::vector< int > score;

  while ((ch = getopt_long(argc, argv, "hvVASNC1M", long_options, &option_index))!=-1) switch(ch) {
    case 0:
      fprintf(stderr, "sanity error, invalid optino to parse, exiting\n");
      exit(-1);
      break;
    case 'S':
      print_filename_flag = 1;
      break;
    case 'N':
      print_index_flag = 1;
      break;
    case 'v':
      print_version();
      exit(0);
      break;
    case 'V':
      opt.verbose=1;
      break;
    case 'A':
      opt.show_all_score=1;
      break;
    case 'C':
      base_concordance=1;
      break;
    case '1':
      coordinate_0ref=0;
      break;
    case 'M':
      case_insensitive=1;
      break;
    default:
    case 'h':
      print_usage();
      exit(0);
      break;
  }

  if (argc>optind) {
    for (i=0; i<(argc-optind-1); i++) {
      ref_fns.push_back( argv[optind+i] );
    }
    ifn = argv[optind+i];
  }

  if (ref_fns.size()==0) {
    printf("Provide input sequences to compare\n");
    print_usage();
    exit(0);
  }

  if (base_concordance) {

    r = do_base_concordance(ifn, ref_fns, &max_score, &max_idx, coordinate_0ref, case_insensitive);

    if (r<0) {
      fprintf(stderr, "ERROR: do_base_concordance got %i, exiting\n", r);
      exit(r);
    }

    if (print_filename_flag) {
      printf("%s\n", ref_fns[max_idx].c_str());
    }
    else if (print_index_flag) {
      printf("%i\n", max_idx);
    }
    else {
      printf("score: %f\nidx: %i\nname: %s\n",
          (float)max_score,
          max_idx,
          ref_fns[max_idx].c_str());
    }

    exit(0);
  }

  //DEBUG
  //
  if (loc_debug) {
    for (i=0; i<ref_fns.size(); i++) {
      printf("%s\n", ref_fns[i].c_str());
    }
    printf("ifn: %s\n", ifn.c_str());
  }

  for (i=0; i<ref_fns.size(); i++) {
    fp = fopen(ref_fns[i].c_str(), "r");
    if (!fp) {
      perror(ref_fns[i].c_str());
      exit(-1);
    }

    tseq.clear();
    while (!feof(fp)) {
      ch = fgetc(fp);
      if (ch==EOF)  { continue; }
      if (ch=='\n') { continue; }
      tseq += (char)ch;
    }
    ref_seq.push_back(tseq);

    fclose(fp);
  }

  fp = fopen(ifn.c_str(), "r");
  if (!fp) {
    perror(ifn.c_str());
    exit(-1);
  }

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF)  { continue; }
    if (ch=='\n') { continue; }
    seq += (char)ch;
  }
  fclose(fp);
  if (loc_debug) {
    for (i=0; i<ref_fns.size(); i++) {
      printf("seq(%s): %s\n",
          ref_fns[i].c_str(),
          ref_seq[i].c_str());
    }
    printf("inp: %s\n", seq.c_str());
  }

  min_score = -1;
  min_idx = 0;
  for (i=0; i<ref_fns.size(); i++) {
    k = asm_ukk_score( (char *)ref_seq[i].c_str(), (char *)seq.c_str() );
    score.push_back(k);

    if ((i==0) || (k<min_score)) {
      min_score = k;
      min_idx = i;
    }
  }

  if (print_filename_flag) {
    printf("%s\n", ref_fns[min_idx].c_str());
  }
  else if (print_index_flag) {
    printf("%i\n", min_idx);
  }
  else {
    printf("min_score: %i\nmin_idx: %i\nname: %s\n",
        min_score,
        min_idx,
        ref_fns[min_idx].c_str());
  }

}
