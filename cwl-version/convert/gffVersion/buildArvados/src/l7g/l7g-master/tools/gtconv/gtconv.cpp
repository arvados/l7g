#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdint.h>
#include <getopt.h>

#include <errno.h>

#include <vector>
#include <string>

#include "pasta.hpp"

#define GTCONV_VERSION "0.1.0"

int _lc(int ch) {
  if ((ch >= 'A') && (ch<='Z')) { return ch + 'a' - 'A'; }
  return ch;
}

int _uc(int ch) {
  if ((ch >= 'a') && (ch<='z')) { return ch - ('a' - 'A'); }
  return ch;
}

typedef struct tile_assembly_type {
  std::vector< std::vector< int > > tile_step, start, len, span;
  std::vector< std::string > name;
} tile_assembly_t;

typedef struct tile_assembly_path_type {
  std::vector< int > tile_step, start, len, span;
  std::string name;
} tile_assembly_path_t;

typedef struct tile_info_type {
  uint32_t tileid;
  int tile_path, tile_step, tile_var;
  std::string tileid_str;

  int span;

  std::string chrom;
  int pos0ref;

} tile_info_t;

typedef struct gtlist_type {
  std::vector< char > allele[2];
  std::vector< std::string > chrom;
  std::vector< int > pos0ref;
  std::vector< std::string > snp_id;
} gtlist_t;

typedef struct gt_tile_info_type {
  int tile_path;
  int allele_count;
  std::vector< std::string > pasta_seq[2];
  std::vector< std::string > tile_seq[2];
  std::vector< int > tile_step;
  std::vector< int > tile_span[2];
  std::vector< std::vector< int > > hiq[2];
  std::vector< int > anchor_tile_step[2];

  std::vector< std::string > chrom;
  std::vector< int > pos0ref_start;
} gt_tile_info_t;

static struct option long_options[] = {
  {"help",                no_argument,        NULL, 'h'},
  {"version",             no_argument,        NULL, 'v'},
  {"verbose",             no_argument,        NULL, 'V'},
  {"input",               required_argument,  NULL, 'i'},
  {"output",              required_argument,  NULL, 'o'},
  {"ref-sequence",        required_argument,  NULL, 'r'},
  {"assembly",            required_argument,  NULL, 'a'},
  {"ref-start0",          required_argument,  NULL, 's'},
  {"ref-start1",          required_argument,  NULL, 'S'},
  {"chrom",               required_argument,  NULL, 'c'},
  {"tile-path",           required_argument,  NULL, 'p'},
  {"allele-count",        required_argument,  NULL, 'A'},
  {0,0,0,0}
};

void show_version() {
  printf("gtconv version %s\n", GTCONV_VERSION);
}

void show_help() {
  printf("\n");
  printf("usage:\n  ngtconv [-h] [-v] [-V] [-i inp] [-o out] [-s s0] [-S s1] [-c chrom] [-p tilepath] -r ref -a assembly [ifn]\n");
  printf("\n");
  printf("  -r ref          reference stream\n");
  printf("  -a assembly     tile assembly stream\n");
  printf("  [-i ifn]        input file\n");
  printf("  [-o ofn]        output (default stdout)\n");
  printf("  [-s s0]         reference stream start, 0 reference (default 0)\n");
  printf("  [-S s1]         reference stream start, 1 reference (default 1)\n");
  printf("  [-c chrom]      chromosome (default 'unk')\n");
  printf("  [-p tilepath]   tile path\n");
  printf("  [-h]            help (this screen)\n");
  printf("  [-v]            version\n");
  printf("  [-V]            verbose\n");
  printf("\n");
}

void print_assembly(tile_assembly_t &ta) {
  int i, j, k, n;
  n = (int)ta.name.size();
  for (i=0; i<n; i++) {
    printf(">%s\n", ta.name[i].c_str());
    for (j=0; j<ta.start[i].size(); j++) {
      printf("%04x\t%10i\n", ta.tile_step[i][j] + ta.span[i][j] - 1, ta.start[i][j] + ta.len[i][j]);
    }
  }

}

void print_tile_assembly_path(tile_assembly_path_t &tap) {
  int i, j, k, n;
  n = (int)tap.tile_step.size();
  printf(">%s\n", tap.name.c_str());
  for (i=0; i<n; i++) {
    printf("%04x\t%10i\n", tap.tile_step[i] + tap.span[i] - 1, tap.start[i] + tap.len[i]);
  }

  for (i=0; i<n; i++) {
    printf("#%04x+%x %i+%i\n", tap.tile_step[i], tap.span[i], tap.start[i] , tap.len[i]);
  }
}

int split_tok_ch(std::vector< std::string > &tok_v, std::string &line, char field) {
  int i, j, k, n;
  std::string buf;

  n = (int)line.size();
  tok_v.clear();

  for (i=0; i<n; i++) {
    if (line[i]==field) {
      if (buf.size()>0) { tok_v.push_back(buf); }
      buf.clear();
      continue;
    }
    buf += line[i];
  }
  if (buf.size()>0) { tok_v.push_back(buf); }

  return 0;
}

int split_tok(std::vector< std::string > &tok_v, std::string &line) {
  int i, j, k, n;
  std::string buf;

  n = (int)line.size();
  tok_v.clear();

  for (i=0; i<n; i++) {
    if ((line[i]==' ') || (line[i] == '\t') || (line[i]=='\n') || (line[i]=='\r')) {
      if (buf.size()>0) { tok_v.push_back(buf); }
      buf.clear();
      continue;
    }
    buf += line[i];
  }
  if (buf.size()>0) { tok_v.push_back(buf); }

  return 0;
}

void print_gt(gtlist_t &gt) {
  int i, j, k, n;

  n = (int)gt.pos0ref.size();

  for (i=0; i<n; i++) {
    printf("%s,%s,%i,%c,%c\n",
        gt.snp_id[i].c_str(),
        gt.chrom[i].c_str(),
        gt.pos0ref[i],
        gt.allele[0][i],
        gt.allele[1][i]);
  }

}

int read_genotype(FILE *fp, gtlist_t &gt) {
  int i, j, k;
  char ch;
  std::string line;
  std::vector< std::string > tok_v;

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {
      if (line.size()==0) { continue; }
      if (line[0]=='#') { line.clear(); continue; }
      split_tok(tok_v, line);

      // 23andme style
      //
      if (tok_v.size()==4) {

        // skip uncalled lines
        //
        if ((tok_v[3] != "--") &&
            (tok_v[3] != "DD") &&
            (tok_v[3] != "II") &&
            (tok_v[3] != "I") &&
            (tok_v[3] != "D") &&
            (tok_v[3] != "DI") &&
            (tok_v[3] != "ID")) {
          gt.snp_id.push_back(tok_v[0]);
          gt.chrom.push_back(tok_v[1]);
          gt.pos0ref.push_back( atoi(tok_v[2].c_str()) - 1 );
          if (tok_v[3].size()==1) {
            gt.allele[0].push_back( tok_v[3][0] );
            //gt.allele[1].push_back( '.' );
            gt.allele[1].push_back( 'n' );
          }
          else if (tok_v[2].size()>=2) {
            gt.allele[0].push_back( tok_v[3][0] );
            gt.allele[1].push_back( tok_v[3][1] );
          }
          else {
            gt.allele[0].push_back( 'n' );
            gt.allele[1].push_back( 'n' );
          }

        }

      }

      //Ancestry style
      if (tok_v.size()==5) {

        // skip uncalled lines
        //
        if ((tok_v[3] != "--") || (tok_v[3] != "-")) {
          gt.snp_id.push_back(tok_v[0]);
          gt.chrom.push_back(tok_v[1]);
          gt.pos0ref.push_back( atoi(tok_v[2].c_str()) - 1 );

          if (tok_v[3].size()==1) { gt.allele[0].push_back( tok_v[3][0] ); }
          else { gt.allele[0].push_back( 'n' ); }

          if (tok_v[4].size()==1) { gt.allele[1].push_back( tok_v[4][0] ); }
          else { gt.allele[1].push_back( 'n' ); }
        }

      }

      line.clear();

    }

    line += (char)ch;
  }

  if ((line.size()>0) && (line[0]!='#')) {
    split_tok(tok_v, line);

    // 23andme style
    //
    if (tok_v.size()==4) {

      // skip uncalled lines
      //
      if (tok_v[3] != "--") {
        gt.snp_id.push_back(tok_v[0]);
        gt.pos0ref.push_back( atoi(tok_v[1].c_str()) - 1 );
        if (tok_v[2].size()==1) {
          gt.allele[0].push_back( tok_v[2][0] );
          gt.allele[1].push_back( 'n' );
        }
        else if (tok_v[2].size()>=2) {
          gt.allele[0].push_back( tok_v[2][0] );
          gt.allele[1].push_back( tok_v[2][1] );
        }
        else {
          gt.allele[0].push_back( 'n' );
          gt.allele[1].push_back( 'n' );
        }

      }

    }

    //Ancestry style
    if (tok_v.size()==5) {

      // skip uncalled lines
      //
      if ((tok_v[3] != "--") || (tok_v[3] != "-")) {
        gt.snp_id.push_back(tok_v[0]);
        gt.pos0ref.push_back( atoi(tok_v[1].c_str()) - 1 );

        if (tok_v[2].size()==1) { gt.allele[0].push_back( tok_v[2][0] ); }
        else { gt.allele[0].push_back( 'n' ); }

        if (tok_v[3].size()==1) { gt.allele[1].push_back( tok_v[3][0] ); }
        else { gt.allele[1].push_back( 'n' ); }
      }

    }


  }

}

int read_tile_assembly_path(FILE *fp, tile_assembly_path_t &tap, int start_0ref) {
  std::string buf, t;
  char ch;
  int i, j, k;
  int prev_tile_step=-1;
  int read_state = 0;
  int x, y;
  std::string prev_chrom, cur_chrom;
  std::vector< std::string > tok_v;

  int tile_step=0, cur_tile_start=0;
  int ref_end_inc_0ref;
  int start_back_offset=0;

  long int li;

  tap.tile_step.clear();
  tap.start.clear();
  tap.len.clear();
  tap.span.clear();
  tap.name.clear();

  cur_tile_start = start_0ref;

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {
      if (buf.size()==0) { continue; }

      if (buf[0]=='>') {

        split_tok_ch(tok_v, buf, ':');
        if (tok_v.size()>=2) { cur_chrom = tok_v[1]; }

        if (prev_chrom != cur_chrom) { cur_tile_start = 0; }
        else { cur_tile_start += 24; }
        prev_chrom = cur_chrom;

        t.clear();
        for (i=1; i<buf.size(); i++) { t += buf[i]; }

        tap.name = t;

        prev_tile_step = -1;
        start_back_offset=0;

      } else {

        t.clear();
        for (i=0; (i<buf.size()) && (buf[i]!='\t') && (buf[i]!=' '); i++) { t.push_back(buf[i]); }
        li = strtol(t.c_str(), NULL, 16);
        tile_step = (int)li;

        tap.tile_step.push_back(prev_tile_step+1);

        t.clear();
        for (; i<buf.size(); i++) { t += buf[i]; }
        ref_end_inc_0ref = atoi(t.c_str());

        tap.start.push_back(cur_tile_start);
        tap.len.push_back(ref_end_inc_0ref - cur_tile_start);
        cur_tile_start = ref_end_inc_0ref - 24;
        tap.span.push_back( tile_step - prev_tile_step );

        prev_tile_step = tile_step;

      }

      buf.clear();
      continue;
    }

    buf += (char)ch;
  }

  return 0;
}

int read_tile_assembly(FILE *fp, tile_assembly_t &ta) {
  std::string buf, t;
  char ch;
  int i, j, k;
  int prev_tile_step=0, tile_path=-1;
  int read_state = 0;
  int x, y;
  std::vector< int > ivec;
  std::string prev_chrom, cur_chrom;
  std::vector< std::string > tok_v;

  int tile_step=0, cur_tile_start=0;
  int ref_end_inc_0ref;
  int start_back_offset=0;

  long int li;

  ivec.clear();

  //ta.end_bound.clear();
  ta.tile_step.clear();
  ta.start.clear();
  ta.len.clear();
  ta.span.clear();
  ta.name.clear();

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {
      if (buf.size()==0) { continue; }

      if (buf[0]=='>') {

        split_tok_ch(tok_v, buf, ':');
        if (tok_v.size()>=2) { cur_chrom = tok_v[1]; }

        if (prev_chrom != cur_chrom) { cur_tile_start = 0; }
        else { cur_tile_start += 24; }
        prev_chrom = cur_chrom;

        ivec.clear();
        t.clear();
        for (i=1; i<buf.size(); i++) { t += buf[i]; }

        ta.name.push_back(t);

        ta.tile_step.push_back(ivec);
        ta.start.push_back(ivec);
        ta.len.push_back(ivec);
        ta.span.push_back(ivec);

        tile_path = (int)ta.tile_step.size()-1;
        prev_tile_step = -1;
        start_back_offset=0;

      } else {
        t.clear();
        for (i=0; (i<buf.size()) && (buf[i]!='\t') && (buf[i]!=' '); i++) { t.push_back(buf[i]); }
        li = strtol(t.c_str(), NULL, 16);
        tile_step = (int)li;

        ta.tile_step[tile_path].push_back(prev_tile_step+1);

        t.clear();
        for (; i<buf.size(); i++) { t += buf[i]; }
        ref_end_inc_0ref = atoi(t.c_str());

        ta.start[tile_path].push_back(cur_tile_start);
        ta.len[tile_path].push_back(ref_end_inc_0ref - cur_tile_start);
        cur_tile_start = ref_end_inc_0ref - 24;
        ta.span[tile_path].push_back( tile_step - prev_tile_step );

        prev_tile_step = tile_step;

      }

      buf.clear();
      continue;
    }

    buf += (char)ch;
  }

  return 0;
}

// process_gt_tile_path stores implied sequence tile information into gt_tile_info_t (gtt_info).
//
// ref_fp         - reference stream starting at `ref_pos_start`
// gtt_info       - result
// gt             - input genotype data (only needs data from relevant tile path)
// assembly       - tile assembly for tile path being processed
// ref_pos_start  - 0 reference start of reference stream
// chrom          - chromosome name
// tile_path      - tile path being processed
//
// The idea is to merge the reference stream (ref_fp) with the genotype data (gt)
// and the assembly information (assembly) to create the tile sequence data.
//
// Each reference character is read and the genotype index is advanced to the appropriate
// position.  If the genotype position matches the reference position, the reported
// base is stored along with the 'high quality' data (the position and length of the reported
// genotype value).  If the reference stream hits the assembly end, save the old tile
// if there was any genotype data for it and create a new tile.
//
// At the end, another pass is done to merge any tiles that have non reference variants
// that fall on a tag and extend the tiles to make them spanning if they weren't already.
//
int process_gt_tile_path(FILE *ref_fp,
                         gt_tile_info_t &gtt_info,
                         gtlist_t &gt,
                         tile_assembly_path_t &assembly,
                         int ref_pos_start,
                         std::string &chrom,
                         int tile_path,
                         int allele_count) {
  int i, j, k, ch, n;
  int ii, jj;

  int ref_pos_len = 0;

  int cur_tile_start = 0;
  int cur_tile_len = 0;
  int gt_idx = 0;
  int last_update_pos = 0;
  int ch_pa;

  std::string ref_str, s;
  std::string allele[2];

  int assembly_idx=0;
  int cur_tile_step = 0;

  int debug_flag = 0;

  std::vector< int > gt_hiq_info[2];
  std::vector< int > ivec;

  while (gt_idx<gt.pos0ref.size()) {
    if (gt.chrom[gt_idx]==chrom) { break; }
    gt_idx++;
  }

  if (debug_flag) {
    printf("starting on gt_idx %i (%s %i)\n", gt_idx, gt.chrom[gt_idx].c_str(), gt.pos0ref[gt_idx]);
  }

  cur_tile_start = assembly.start[0];
  cur_tile_len = assembly.len[0];

  gtt_info.tile_path = tile_path;

  while (!feof(ref_fp)) {
    ch = fgetc(ref_fp);
    if (ch==EOF) { continue; }
    if ((ch=='\n') || (ch==' ') || (ch=='\r')) { continue; }

    // advance gt_idx until we reach our current reference position
    //
    while (gt_idx < (int)gt.pos0ref.size()) {
      if (gt.pos0ref[gt_idx] >= (ref_pos_start+ref_pos_len)) { break; }

      if (debug_flag) {
      printf("    gt[%i] %s:%i (%c%c)\n",
          gt_idx,
          gt.chrom[gt_idx].c_str(),
          gt.pos0ref[gt_idx],
          gt.allele[0][gt_idx],
          gt.allele[1][gt_idx]);
      }

      gt_idx++;
    }

    // Add base to allele
    //

    /*
    if (gt_idx >= gt.pos0ref.size()) {
      printf("#ERROR gt_idx %i >= pos0ref.size() %i\n",
          gt_idx,
          (int)gt.pos0ref.size());
    }

    printf("#debug: gt_idx %i\n", gt_idx);
    printf("#debug: pos0ref %i\n", gt.pos0ref[gt_idx]);
    printf("#debug: ref_pos_start %i\n", ref_pos_start);
    printf("#debug: ref_pos_len %i\n", ref_pos_len);
    */

    if ((gt_idx < gt.pos0ref.size()) &&
        (gt.pos0ref[gt_idx] == (ref_pos_start+ref_pos_len))) {

      ch_pa = pasta_convert(ch, _lc(gt.allele[0][gt_idx]));
      if (ch_pa==0) {
        fprintf(stderr, "INVALID PASTA CONVERSION %c (%i) -> %c (%i) at %i\n",
            ch, ch,
            gt.allele[0][gt_idx], gt.allele[0][gt_idx],
           ref_pos_start+ref_pos_len);
        exit(-1);
      }
      allele[0] += (char)ch_pa;

      if (gt.allele[1][gt_idx]!='.') {
        ch_pa = pasta_convert(ch, _lc(gt.allele[1][gt_idx]));
        if (ch_pa==0) { perror("INVALID PASTA CONVERSION\n"); exit(-1); }
        allele[1] += (char)ch_pa;
      }
      else {
        allele[1] += '.';
      }

      last_update_pos = ref_pos_start + ref_pos_len;

      if (debug_flag) {
      printf("++ gt %c%c (%s:%i)\n",
          gt.allele[0][gt_idx],
          gt.allele[1][gt_idx],
          gt.chrom[gt_idx].c_str(),
          gt.pos0ref[gt_idx]);
      }

      gt_hiq_info[0].push_back(ref_pos_len); gt_hiq_info[0].push_back(1);
      gt_hiq_info[1].push_back(ref_pos_len); gt_hiq_info[1].push_back(1);

    }
    else {
      allele[0] += (char)ch;
      allele[1] += (char)ch;
    }

    // update state
    //
    ref_pos_len++;
    ref_str += (char)ch;

    // Advance to next tile
    //
    if (ref_pos_len == cur_tile_len) {

      // If we have genotype information on a tile, emit information
      //
      if (last_update_pos >= ref_pos_start) {

        if (debug_flag) {
          printf("++ gtt_info tile_step %i (assembly_idx %i)\n", assembly.tile_step[assembly_idx], assembly_idx);
        }

        gtt_info.pasta_seq[0].push_back(allele[0]);
        gtt_info.pasta_seq[1].push_back(allele[1]);
        gtt_info.tile_step.push_back(assembly.tile_step[assembly_idx]);
        gtt_info.tile_span[0].push_back(assembly.span[assembly_idx]);
        gtt_info.tile_span[1].push_back(assembly.span[assembly_idx]);

        gtt_info.hiq[0].push_back( gt_hiq_info[0] );
        gtt_info.hiq[1].push_back( gt_hiq_info[1] );

        gtt_info.chrom.push_back(chrom);
        gtt_info.pos0ref_start.push_back(assembly.start[assembly_idx]);
      }

      assembly_idx++;

      if (debug_flag) {
        if (assembly_idx < assembly.tile_step.size()) {
          printf("# assembly advanced to [%i] %x.%x+%x:%i+%i (%s)\n",
              assembly_idx,
              tile_path,
              assembly.tile_step[assembly_idx],
              assembly.span[assembly_idx],
              assembly.start[assembly_idx],
              assembly.len[assembly_idx],
              assembly.name.c_str());
        }
      }

      if (assembly_idx < assembly.len.size()) {
        cur_tile_len = assembly.len[assembly_idx];
        cur_tile_start = assembly.start[assembly_idx];
      }


      // reset state
      //
      ivec.clear();
      n = allele[0].size();
      int pos_threshold = n-24;
      for (ii=0; ii<gt_hiq_info[0].size(); ii+=2) {
        if (gt_hiq_info[0][ii] >= pos_threshold) {
          ivec.push_back(gt_hiq_info[0][ii] - pos_threshold);
          ivec.push_back(gt_hiq_info[0][ii+1]);
        }
      }
      gt_hiq_info[0].clear();
      gt_hiq_info[0] = ivec;

      ivec.clear();
      n = allele[1].size();
      pos_threshold = n-24;
      for (ii=0; ii<gt_hiq_info[1].size(); ii+=2) {
        if (gt_hiq_info[1][ii] >= pos_threshold) {
          ivec.push_back(gt_hiq_info[1][ii] - pos_threshold);
          ivec.push_back(gt_hiq_info[1][ii+1]);
        }
      }
      gt_hiq_info[1].clear();
      gt_hiq_info[1] = ivec;

      ivec.clear();

      n = (int)ref_str.size();
      s.clear();
      for (i=0; i<24; i++) { s += ref_str[n-24+i]; }
      ref_str.clear();
      ref_str += s;

      n = (int)allele[0].size();
      s.clear();
      for (i=0; i<24; i++) { s += allele[0][n-24+i]; }
      allele[0].clear();
      allele[0] += s;

      n = (int)allele[1].size();
      s.clear();
      for (i=0; i<24; i++) { s += allele[1][n-24+i]; }
      allele[1].clear();
      allele[1] += s;

      ref_pos_start += ref_pos_len - 24;
      ref_pos_len = 24;

    }

  }

  // If we have (non-ref) variants that fall on tags,
  // we need to merge the tiles together.
  // As a final post-processing step, run through all
  // tiles/sequences and merge tiles we've marked
  // above as having non-ref variants on tags.
  //

  int anchor_tile_idx[2], anchor_tile_step[2];
  int iallele=0;

  anchor_tile_step[0]=0;
  anchor_tile_step[1]=0;

  anchor_tile_idx[0] = 0;
  anchor_tile_idx[1] = 0;

  for (iallele=0; iallele<allele_count; iallele++) {
  int merge_needed = 0;

  for (i=0; i<gtt_info.pasta_seq[iallele].size(); i++) {

    merge_needed = 0;

    // check to see if we need to merge tiles
    // if we're at the beginning of the path, having
    // a non-ref variant on the tag doesn't imply
    // a merge as there's no real tag there and
    // nothing to merge with.
    //
    if (gtt_info.tile_step[i] > 0) {
      n = (int)gtt_info.pasta_seq[iallele][i].size();
      for (k=0; k<24; k++) {
        //ch_pa = gtt_info.pasta_seq[iallele][i][k];
        ch_pa = _lc(gtt_info.pasta_seq[iallele][i][k]);
        if ((ch_pa != 'n') &&
            (ch_pa != 'a') &&
            (ch_pa != 'c') &&
            (ch_pa != 'g') &&
            (ch_pa != 't')) {

          if (debug_flag) {
            printf("********** merge needed (idx %i, tile_step %i)\n",
                i, gtt_info.tile_step[i]);
          }

          merge_needed = 1;
          break;
        }
      }
    }

    // If we don't need to merge, update our anchor step
    //
    if (!merge_needed) {
      anchor_tile_step[iallele] = (int)gtt_info.tile_step[i];
      anchor_tile_idx[iallele] = i;
    }

    // Add our anchor step information to see if we need
    // to skip this sequence below.
    //
    gtt_info.anchor_tile_step[iallele].push_back( anchor_tile_step[iallele] );

    if (debug_flag) {
    printf("# tile step %i (%x), anchor step %i (%x)\n",
        gtt_info.tile_step[i],
        gtt_info.tile_step[i],
        gtt_info.anchor_tile_step[iallele][i],
        gtt_info.anchor_tile_step[iallele][i]);
    }

    s.clear();
    for (k=0; k<gtt_info.pasta_seq[iallele][i].size(); k++) {
      ch_pa = pasta2seq(gtt_info.pasta_seq[iallele][i][k]);
      if (ch_pa > 0) { s += (char)ch_pa; }
    }

    gtt_info.tile_seq[iallele].push_back(s);

    // Add the sequence (sans beginning tag) to the anchor sequence
    //
    if (merge_needed) {

      n = (int)gtt_info.tile_seq[iallele].size();
      int m = (int)gtt_info.tile_seq[iallele][n-1].size();
      for (k=24; k<m; k++) {
        gtt_info.tile_seq[iallele][anchor_tile_idx[iallele]] += gtt_info.tile_seq[iallele][n-1][k];
      }

      //n = (int)gtt_info.pasta_seq[iallele].size();
      //m = (int)gtt_info.pasta_seq[iallele][n-1].size();

      /*
      printf("### pasts n %i, m %i, a %i, step %i (%x)... merging in step %i (%x)\n",
          n, m,
          iallele,
          (int)gtt_info.pasta_seq[iallele].size(),
          (int)gtt_info.pasta_seq[iallele].size(),
          n-1,
          n-1
          );
          */

      for (k=24; k<m; k++) {
        gtt_info.pasta_seq[iallele][anchor_tile_idx[iallele]] += gtt_info.pasta_seq[iallele][n-1][k];
      }

      gtt_info.tile_span[iallele][anchor_tile_idx[iallele]] += gtt_info.tile_span[iallele][n-1];

      // copy over 'hiq' values.
      // even entries are position which we need to renormalize to the beginning of this tile
      // odd entreis are length which can be transferred over without issue
      //
      for (k=0; k<gtt_info.hiq[iallele][n-1].size(); k+=2) {

        int v = gtt_info.hiq[iallele][n-1][k];

        // any variants that already fall on a tag should already be in the 'hiq' list
        // of the anchor tile so we can skip over them.
        //
        if (v<24) { continue; }
        v += gtt_info.pos0ref_start[n-1] - gtt_info.pos0ref_start[anchor_tile_idx[iallele]];

        gtt_info.hiq[iallele][anchor_tile_idx[iallele]].push_back(v);
        gtt_info.hiq[iallele][anchor_tile_idx[iallele]].push_back( gtt_info.hiq[iallele][n-1][k+1] );

      }
    }

  }
  }

  return 0;

}

// output information on each tile.
// output format is:
//
// <tileid>,0,<seq>
// <tileid>,1,<chrom>,<start_tile_ref0pos>,[ <start_offset_0ref> <len> .. ]
//
// With the different alleles indicated in the appropriate variant.
// The print order should not be assumed and a 'sort' should be issued
// afterwards if a consistent sort order is desired.
//

void print_gt_tile_info(gt_tile_info_t &gtt_info, int allele_count) {
  int a, i, j, k;
  int ii, jj;
  int skip_tile = 0;
  int debug_flag = 0;
  int tile_path;

  tile_path = gtt_info.tile_path;

  for (a=0; a<allele_count; a++) {

    for (i=0; i<gtt_info.tile_seq[a].size(); i++) {

      skip_tile = 0;
      if (gtt_info.anchor_tile_step[a][i] != gtt_info.tile_step[i]) {

        if (debug_flag) {
        printf("# idx %i, tile_step %i spanning (anchor %i), skipping\n",
            i, gtt_info.tile_step[i],
            gtt_info.anchor_tile_step[a][i]);
        }

        skip_tile=1;
        continue;
      }

      printf("%04x.%02x.%04x.%03x+%x,0,%s\n",
          tile_path, 0, gtt_info.tile_step[i], a,
          gtt_info.tile_span[a][i],
          gtt_info.tile_seq[a][i].c_str());

      printf("%04x.%02x.%04x.%03x+%x,1,",
          tile_path, 0, gtt_info.tile_step[i], a,
          gtt_info.tile_span[a][i]);
      printf("%s,%i,", gtt_info.chrom[i].c_str(), gtt_info.pos0ref_start[i]);

      printf("[");
      for (ii=0; ii<gtt_info.hiq[a][i].size(); ii+=2) {
        printf(" %i %i",
            gtt_info.hiq[a][i][ii],
            gtt_info.hiq[a][i][ii+1]);
      }
      printf("]\n");

    }

  }

}

void print_gt_pasta_tile_info(gt_tile_info_t &gtt_info, int allele_count) {
  int a, i, j, k;
  int ii, jj;
  int skip_tile = 0;
  int debug_flag = 0;
  int tile_path;

  tile_path = gtt_info.tile_path;

  for (a=0; a<allele_count; a++) {

    for (i=0; i<gtt_info.pasta_seq[a].size(); i++) {

      skip_tile = 0;
      if (gtt_info.anchor_tile_step[a][i] != gtt_info.tile_step[i]) {

        if (debug_flag) {
        printf("# idx %i, tile_step %i spanning (anchor %i), skipping\n",
            i, gtt_info.tile_step[i],
            gtt_info.anchor_tile_step[a][i]);
        }

        skip_tile=1;
        continue;
      }

      printf("%04x.%02x.%04x.%03x+%x,0,%s\n",
          tile_path, 0, gtt_info.tile_step[i], a,
          gtt_info.tile_span[a][i],
          gtt_info.pasta_seq[a][i].c_str());

      printf("%04x.%02x.%04x.%03x+%x,1,",
          tile_path, 0, gtt_info.tile_step[i], a,
          gtt_info.tile_span[a][i]);
      printf("%s,%i,", gtt_info.chrom[i].c_str(), gtt_info.pos0ref_start[i]);

      printf("[");
      for (ii=0; ii<gtt_info.hiq[a][i].size(); ii+=2) {
        printf(" %i %i",
            gtt_info.hiq[a][i][ii],
            gtt_info.hiq[a][i][ii+1]);
      }
      printf("]\n");

    }

  }

}

int main(int argc, char **argv) {
  int i, j, k;
  int opt;
  int ret=0;
  int idx;
  char buf[1024];

  FILE *ifp=stdin,
       *ofp=stdout,
       *ref_fp=NULL,
       *assembly_fp=NULL;

  int option_index=0;
  int verbose_flag = 0;
  tile_assembly_t assembly;
  tile_assembly_path_t assembly_path;
  gtlist_t gt;

  int tile_path = -1;

  int ref_start0=0;
  std::string chrom = "unk";
  int allele_count=2;

  gt_tile_info_t gtt_info;

  std::string ifn="-",
              ofn="-",
              ref_fn,
              assembly_fn;

  while ((opt = getopt_long(argc, argv, "hvVi:o:r:a:s:S:c:p:A:", long_options, &option_index))!=-1) switch (opt) {
    case 0:
      fprintf(stderr, "invalid option, exiting\n");
      exit(1);
      break;
    case 'v': show_version(); exit(0); break;
    case 'V': verbose_flag=1; break;
    case 'i': ifn=optarg; break;
    case 'o': ofn=optarg; break;
    case 'r': ref_fn=optarg; break;
    case 'a': assembly_fn = optarg; break;
    case 's': ref_start0 = atoi(optarg); break;
    case 'S': ref_start0 = atoi(optarg)-1; break;
    case 'c': chrom = optarg; break;
    case 'p': tile_path = atoi(optarg); break;
    case 'A': allele_count=atoi(optarg); break;
    case 'h':
    default: show_help(); exit(0); break;
  }

  if ((ref_fn.size()==0) || (assembly_fn.size()==0)) {
    printf("provide reference and assembly streams\n");
    show_help();
    exit(0);
  }

  if (ifn!="-") {
    if (!(ifp = fopen(ifn.c_str(),"r"))) {
      perror(ifn.c_str());
      exit(2);
    }
  }

  if (ofn!="-") {
    if (!(ofp = fopen(ofn.c_str(), "w"))) {
      perror(ofn.c_str());
      exit(3);
    }
  }

  if (!(ref_fp = fopen(ref_fn.c_str(), "r"))) {
    perror(ref_fn.c_str());
    exit(4);
  }

  if (!(assembly_fp = fopen(assembly_fn.c_str(), "r"))) {
    perror(assembly_fn.c_str());
    exit(5);
  }

  if (verbose_flag) { printf("# read assembly.."); fflush(stdout); }

  read_tile_assembly_path(assembly_fp, assembly_path, ref_start0);

  if (verbose_flag) { printf("ok\n");  }

  if (verbose_flag) { printf("# reading genotype..."); fflush(stdout); }

  read_genotype(ifp, gt);

  if (verbose_flag) { printf("ok\n"); }

  process_gt_tile_path(ref_fp, gtt_info, gt, assembly_path, ref_start0, chrom, tile_path, allele_count);

  //print_gt_tile_info(gtt_info, allele_count);
  print_gt_pasta_tile_info(gtt_info, allele_count);

  if (ifp!=stdin) { fclose(ifp); }
  if (ofp!=stdout) { fclose(ofp); }
  if (ref_fp) { fclose(ref_fp); }
  if (assembly_fp) { fclose(assembly_fp); }

}

