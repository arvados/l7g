#include "sglf.hpp"
#include "tileband-hash.hpp"

int band_md5_hash(std::vector< std::string > &digest_str, std::vector< band_info_t > &band_v, sglf_t &sglf, std::vector< int > &tilepath_list) {
  int i, j, k, n, m, a;

  int allele=0, tilepath=0, tilestep=0, tilevar=0, span_len=0;
  int noc_count=0, noc_start, noc_len, pos;

  int start_idx, idx, del_idx, ds=0;

  MD5_CTX md5_ctx[2];
  unsigned char digest[MD5_DIGEST_LENGTH];
  char cbuf[128], *chp;

  std::string hash, hash_mask, seq, seq_mask, m5_s;

  k = (int)(band_v.size()) % (int)(tilepath_list.size());
  if (k!=0) { return -1; }

  del_idx = (int)tilepath_list.size();
  m5_s.clear();
  cbuf[0]='\0';

  for (ds=0, start_idx=0; start_idx<band_v.size(); start_idx += del_idx, ds++) {

    MD5_Init(&(md5_ctx[0]));
    MD5_Init(&(md5_ctx[1]));

    for (idx=start_idx; idx < (start_idx+del_idx); idx++) {

      if (tilepath >= (int)sglf.seq.size()) { return -2; }

      tilepath = tilepath_list[idx-start_idx];

      for (allele=0; allele<2; allele++) {

        tilestep = 0;
        while (tilestep < band_v[idx].band[allele].size()) {

          if (tilestep >= sglf.seq[tilepath].size()) { return -3; }

          span_len=1;
          while ( ((tilestep + span_len) < band_v[idx].band[allele].size()) &&
                  (band_v[idx].band[allele][tilestep+span_len]==-1) ) {
            span_len++;
          }

          tilevar = band_v[idx].band[allele][tilestep];

          if (tilevar >= sglf.seq[tilepath][tilestep].size()) { return -4; }

          seq = sglf.seq[tilepath][tilestep][tilevar];

          noc_count=0;
          for (i=0; i<band_v[idx].noc[allele][tilestep].size(); i+=2) {

            noc_start = band_v[idx].noc[allele][tilestep][i];
            noc_len = band_v[idx].noc[allele][tilestep][i+1];

            for (pos=noc_start; pos<(noc_start + noc_len); pos++) {
              seq[pos] = 'n';
            }

            noc_count += noc_len;

          }

          if (tilestep==0) {
            MD5_Update(&(md5_ctx[allele]), (const void *)(seq.c_str()), (unsigned long)seq.size());
          } else {

            if (seq.size()>24) {
              MD5_Update(&(md5_ctx[allele]), (const void *)(seq.c_str()+24), (unsigned long)(seq.size()-24));
            }

          }

          tilestep+=span_len;

        }

      }

    }

    m5_s.clear();
    MD5_Final(digest, &(md5_ctx[0]));
    for (i=0; i<MD5_DIGEST_LENGTH; i++) {
      sprintf(cbuf, "%02x", (unsigned int)digest[i]);
      m5_s += cbuf;
    }

    m5_s += " ";

    MD5_Final(digest, &(md5_ctx[1]));
    for (i=0; i<MD5_DIGEST_LENGTH; i++) {
      sprintf(cbuf, "%02x", (unsigned int)digest[i]);
      m5_s += cbuf;
    }

    digest_str.push_back(m5_s);

  }

  return 0;
}

//--

int read_bands(FILE *ifp, std::vector< band_info_t > &band_info_v) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;
  std::vector< int > noc_vec;

  std::string buf;

  int read_state = 0;
  int bracket_count=0;
  int cur_val=-3;

  band_info_t band_info;

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }

    char_no++;
    if (ch=='\n') {
      line_no++;

      switch(read_state) {
        case 0:
          read_state++;
          break;
        case 1:
          read_state++;
          break;
        case 2:
          read_state++;
          break;
        case 3:
          read_state=0;
          bracket_count=0;
          buf.clear();
          band_info_v.push_back(band_info);
          band_info.band[0].clear();
          band_info.band[1].clear();
          band_info.noc[0].clear();
          band_info.noc[1].clear();
          break;
        default:
          return -1;
      }
      continue;
    }

    if (ch==' ') {
      if (buf.size()>0) {
        cur_val = atoi(buf.c_str());

        if (read_state < 2) {
          band_info.band[read_state].push_back(cur_val);
        }

        else {
          noc_vec.push_back(cur_val);
        }

      }
      buf.clear();
      continue;
    }

    if (ch=='[') { bracket_count++; continue; }
    if (ch==']') {
      bracket_count--;

      // Tile variant bands still
      //
      if (read_state<2) {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());

          if (read_state < 2) {
            band_info.band[read_state].push_back(cur_val);
          }
          buf.clear();
        }

      }

      // nocall information
      //
      else {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());
          noc_vec.push_back(cur_val);
          buf.clear();
        }

        if (bracket_count==1) {
          band_info.noc[read_state-2].push_back(noc_vec);
          noc_vec.clear();
        }
      }

      continue;
    }

    buf += (char)ch;

  }

  return 0;
}


int read_band(FILE *ifp, band_info_t &band_info) {
  int i, j,k ;
  int line_no=0, char_no=0;
  int ch;
  std::vector< int > noc_vec;

  std::string buf;

  int read_state = 0;
  int bracket_count=0;
  int cur_val=-3;

  while (!feof(ifp)) {
    ch = fgetc(ifp);
    if (ch==EOF) { continue; }
    char_no++;
    if (ch=='\n') {
      line_no++;

      switch(read_state) {
        case 0:
          break;
        case 1:
          break;
        case 2:
          break;
        case 3:
          break;
        default:
          return -1;
      }
      read_state++;
      continue;
    }

    if (ch==' ') {
      if (buf.size()>0) {
        cur_val = atoi(buf.c_str());

        if (read_state < 2) {
          band_info.band[read_state].push_back(cur_val);
        }

        else {
          noc_vec.push_back(cur_val);
        }

      }
      buf.clear();
      continue;
    }

    if (ch=='[') { bracket_count++; continue; }
    if (ch==']') {
      bracket_count--;

      // Tile variant bands still
      //
      if (read_state<2) {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());

          if (read_state < 2) {
            band_info.band[read_state].push_back(cur_val);
          }
          buf.clear();
        }

      }

      // nocall information
      //
      else {

        if (buf.size()>0) {
          cur_val = atoi(buf.c_str());
          noc_vec.push_back(cur_val);
          buf.clear();
        }

        if (bracket_count==1) {
          band_info.noc[read_state-2].push_back(noc_vec);
          noc_vec.clear();
        }
      }

      continue;
    }

    buf += (char)ch;

  }

  return 0;
}

void band_print(band_info_t &band_info) {
  int i, j, a;

  for (a=0; a<2; a++) {
    printf("[");
    for (i=0; i<band_info.band[a].size(); i++) {
      printf(" %i", band_info.band[a][i]);
    }
    printf("]\n");
  }

  for (a=0; a<2; a++) {
    printf("[");
    for (i=0; i<band_info.noc[a].size(); i++) {
      printf("[");
      for (j=0; j<band_info.noc[a][i].size(); j++) {
        printf(" %i", band_info.noc[a][i][j]);
      }
      printf(" ]");
    }
    printf("]\n");
  }


}

void print_bands(std::vector< band_info_t > &band_info_v) {
  int ii;
  for (ii=0; ii<band_info_v.size(); ii++) {
    band_print(band_info_v[ii]);
  }
}



