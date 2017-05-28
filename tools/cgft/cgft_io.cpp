#include "cgft.hpp"

static int _write_string(std::string &s, FILE *fp) {
  uint32_t u32;
  size_t sz;

  u32 = (uint32_t)s.size();
  sz = fwrite(&(u32), sizeof(uint32_t), 1, fp);
  if (sz!=1) { return -1; }

  sz = fwrite(&(s[0]), sizeof(char), s.size(), fp);
  if (sz!=s.size()) { return -1; }

  return 0;
}

struct membuf : std::streambuf {
  membuf(char *beg, char *end) { this->setg(beg, beg, end); }
};

static int _read_vlc_vector(sdsl::vlc_vector<> &vlc, size_t sz, FILE *fp) {
  size_t s;
  int ch;
  std::vector<char> b;

  for (s=0; s<sz; s++) {
    ch = fgetc(fp);
    if (ch==EOF) { return -1; }
    b.push_back(ch);
  }

  membuf sbuf(&(b[0]), &(b[0]) + b.size());
  std::istream in(&sbuf);
  vlc.load(in);

  return 0;
}


static int _write_vlc_vector(sdsl::vlc_vector<> &vlc, FILE *fp) {
  std::ostringstream bufstream;
  std::string s;
  size_t sz;

  vlc.serialize(bufstream);
  s = bufstream.str();

  sz = fwrite( s.c_str(), sizeof(char), s.size(), fp);
  if (sz!=s.size()) { return -1; }

  return 0;
}

static int _read_enc_vector(sdsl::enc_vector<> &enc, size_t sz, FILE *fp) {
  size_t s;
  int ch;
  std::vector<char> b;

  for (s=0; s<sz; s++) {
    ch = fgetc(fp);
    if (ch==EOF) { return -1; }
    b.push_back(ch);
  }

  membuf sbuf(&(b[0]), &(b[0]) + b.size());
  std::istream in(&sbuf);
  enc.load(in);

  return 0;
}


static int _write_enc_vector(sdsl::enc_vector<> &enc, FILE *fp) {
  std::ostringstream bufstream;
  std::string s;
  size_t sz;

  enc.serialize(bufstream);
  s = bufstream.str();

  sz = fwrite( s.c_str(), sizeof(char), s.size(), fp);
  if (sz!=s.size()) { return -1; }

  return 0;
}

static uint64_t _calc_path_size(tilepath_t *tilepath) {
  int i, j, k;
  uint32_t u32;
  uint64_t u64;
  uint64_t byte_count = 0, n_8, n_32;

  // TilePath
  //
  byte_count += (uint64_t)sizeof(uint64_t);

  // Name
  //
  byte_count += (uint64_t)(sizeof(uint32_t) + tilepath->Name.size());

  // NTileStep
  // NOverflow
  // NOverflow64
  // ExtraDataSize
  //
  byte_count += (uint64_t)sizeof(uint64_t);
  byte_count += (uint64_t)sizeof(uint64_t);
  byte_count += (uint64_t)sizeof(uint64_t);
  byte_count += (uint64_t)sizeof(uint64_t);

  // Loq
  // Span
  //
  n_8 = ((tilepath->NTileStep+7)/8);
  byte_count += n_8;
  byte_count += n_8;

  // Cache
  //
  n_32 = ((tilepath->NTileStep+31)/32);
  byte_count += n_32;

  byte_count += (uint64_t)(sizeof(uint16_t))*3*(tilepath->NOverflow);
  byte_count += (uint64_t)(sizeof(uint64_t))*3*(tilepath->NOverflow64);

  byte_count += (uint64_t)(sizeof(char)*(tilepath->ExtraDataSize));

  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileStepHom));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileVariantHom));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocSumHom));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocStartHom));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocLenHom));

  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileStepHet));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileVariantHet));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocSumHet));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocStartHet));
  byte_count += (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocLenHet));

  return byte_count;

}

int cgft_write_to_file(cgf_t *cgf, const char *ofn) {
  int i, j, k, n, ii;
  size_t sz;
  uint64_t u64;
  uint32_t u32;
  uint16_t u16;
  FILE *ofp;
  char c;
  int n_8, n_32;
  tilepath_t *tilepath;
  uint64_t byte_count=0;

  unsigned char *xx;

  if ((ofp=fopen(ofn, "w"))==NULL) { return -1; }

  //DEBUG
  //if ((ofp=fopen("dm.cgf3", "w"))==NULL) { return -1; }

  for (i=0; i<8; i++) {
    sz = fwrite(&(cgf->Magic[i]), sizeof(char), 1, ofp);
    if (sz!=1) { return -1; }
    byte_count++;
  }

  k = _write_string(cgf->CGFVersion, ofp);
  if (k<0) { return k; }
  byte_count+=(uint64_t)sizeof(uint32_t) + (uint64_t)cgf->CGFVersion.size();

  k = _write_string(cgf->LibraryVersion, ofp);
  if (k<0) { return k; }
  byte_count+=(uint64_t)sizeof(uint32_t) + (uint64_t)cgf->LibraryVersion.size();


  sz = fwrite(&(cgf->PathCount), sizeof(uint64_t), 1, ofp);
  if (sz!=1) { return -1; }
  byte_count+=(uint64_t)sizeof(uint64_t);

  k = _write_string(cgf->TileMap, ofp);
  if (k<0) { return k; }
  byte_count+=(uint64_t)sizeof(uint32_t) + (uint64_t)cgf->TileMap.size();

  byte_count += (uint64_t)(sizeof(uint64_t))*(cgf->PathCount);

  if (cgf->PathCount>0) {

    cgf->PathStructOffset.clear();
    cgf->PathStructOffset.push_back(byte_count);
    for (ii=1; ii<cgf->Path.size(); ii++) {
      byte_count += _calc_path_size(&(cgf->Path[ii-1]));
      cgf->PathStructOffset.push_back(byte_count);
    }

    sz = fwrite(&(cgf->PathStructOffset[0]), sizeof(uint64_t), (size_t)cgf->PathCount, ofp);
    if (sz!=(size_t)cgf->PathCount) { return -1; }

    for (ii=0; ii<cgf->Path.size(); ii++) {

      tilepath = &(cgf->Path[ii]);

      sz = fwrite(&(tilepath->TilePath), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      u32 = (uint32_t)tilepath->Name.size();
      sz = fwrite(&u32, sizeof(uint32_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->Name[0]), sizeof(char), tilepath->Name.size(), ofp);
      if (sz!=tilepath->Name.size()) { return -1; }

      sz = fwrite(&(tilepath->NTileStep), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->NOverflow), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->NOverflow64), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->ExtraDataSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      n_8 = (int)((tilepath->NTileStep+7)/8);
      sz = fwrite(tilepath->Loq, sizeof(char), (size_t)n_8, ofp);
      if (sz!=(size_t)n_8) { return -1; }

      sz = fwrite(tilepath->Span, sizeof(char), (size_t)n_8, ofp);
      if (sz!=(size_t)n_8) { return -1; }

      n_32 = (int)((tilepath->NTileStep+31)/32);

      sz = fwrite(tilepath->Cache, sizeof(uint64_t), (size_t)n_32, ofp);
      if (sz!=(size_t)n_32) { return -1; }

      sz = fwrite(&(tilepath->Overflow[0]), sizeof(uint16_t), (size_t)(tilepath->NOverflow), ofp);
      if (sz != (size_t)(tilepath->NOverflow)) { return -1; }

      if (tilepath->NOverflow64>0) {
        sz = fwrite(&(tilepath->Overflow64[0]), sizeof(uint64_t), (size_t)(3*(tilepath->NOverflow64)), ofp);
        if (sz != (size_t)(3*tilepath->NOverflow64)) { return -1; }
      }

      if (tilepath->ExtraDataSize>0) {
        sz = fwrite(&(tilepath->ExtraData[0]), sizeof(char), (size_t)(tilepath->ExtraDataSize), ofp);
        if (sz != (size_t)(tilepath->ExtraDataSize)) { return -1; }
      }

      tilepath->LoqTileStepHomSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileStepHom));
      tilepath->LoqTileVariantHomSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileVariantHom));
      tilepath->LoqTileNocSumHomSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocSumHom));
      tilepath->LoqTileNocStartHomSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocStartHom));
      tilepath->LoqTileNocLenHomSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocLenHom));

      tilepath->LoqTileStepHetSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileStepHet));
      tilepath->LoqTileVariantHetSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileVariantHet));
      tilepath->LoqTileNocSumHetSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocSumHet));
      tilepath->LoqTileNocStartHetSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocStartHet));
      tilepath->LoqTileNocLenHetSize = (uint64_t)(sdsl::size_in_bytes(tilepath->LoqTileNocLenHet));

      // size information for hom low quality structures
      //

      sz = fwrite(&(tilepath->LoqTileStepHomSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileVariantHomSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileNocSumHomSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileNocStartHomSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileNocLenHomSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      // size information for het low quality structures
      //

      sz = fwrite(&(tilepath->LoqTileStepHetSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileVariantHetSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileNocSumHetSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileNocStartHetSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      sz = fwrite(&(tilepath->LoqTileNocLenHetSize), sizeof(uint64_t), 1, ofp);
      if (sz!=1) { return -1; }

      // low quality hom structures
      //

      k = _write_enc_vector(tilepath->LoqTileStepHom, ofp);
      if (k<0) { return -1; }

      k = _write_vlc_vector(tilepath->LoqTileVariantHom, ofp);
      if (k<0) { return -1; }

      k = _write_enc_vector(tilepath->LoqTileNocSumHom, ofp);
      if (k<0) { return -1; }

      k = _write_vlc_vector(tilepath->LoqTileNocStartHom, ofp);
      if (k<0) { return -1; }

      k = _write_vlc_vector(tilepath->LoqTileNocLenHom, ofp);
      if (k<0) { return -1; }

      // low quality het structures
      //

      k = _write_enc_vector(tilepath->LoqTileStepHet, ofp);
      if (k<0) { return -1; }

      k = _write_vlc_vector(tilepath->LoqTileVariantHet, ofp);
      if (k<0) { return -1; }

      k = _write_enc_vector(tilepath->LoqTileNocSumHet, ofp);
      if (k<0) { return -1; }

      k = _write_vlc_vector(tilepath->LoqTileNocStartHet, ofp);
      if (k<0) { return -1; }

      k = _write_vlc_vector(tilepath->LoqTileNocLenHet, ofp);
      if (k<0) { return -1; }

    }

  }

  fclose(ofp);

}

int cgft_tilepath_read(tilepath_t *tilepath, FILE *fp) {
  int i, j, k, n;
  int ch;
  uint64_t u64;
  uint32_t u32;
  unsigned char ub[32];
  size_t n_8, n_32, sz;

  tilepath->Loq = NULL;
  tilepath->Span = NULL;
  tilepath->Cache = NULL;
  tilepath->Overflow = NULL;
  tilepath->Overflow64 = NULL;
  tilepath->ExtraData = NULL;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->TilePath = u64;

  sz = fread(&u32, sizeof(uint32_t), 1, fp);
  if (sz!=1) { return -1; }

  tilepath->Name.clear();
  tilepath->Name.reserve(u32);
  for (i=0; i<u32; i++) {
    ch = fgetc(fp);
    if (ch==EOF) { return -1; }
    tilepath->Name += ch;
  }

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->NTileStep = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->NOverflow = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->NOverflow64 = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->ExtraDataSize = u64;

  n_8 = (size_t)((tilepath->NTileStep+7)/8);
  n_32 = (size_t)((tilepath->NTileStep+31)/32);

  tilepath->Loq = new unsigned char[n_8];
  tilepath->Span = new unsigned char[n_8];
  tilepath->Cache = new uint64_t[n_32];

  if (tilepath->NOverflow > 0){
    tilepath->Overflow = new uint16_t[tilepath->NOverflow];
  }

  if (tilepath->NOverflow64 > 0) {
    tilepath->Overflow64 = new uint64_t[tilepath->NOverflow64];
  }

  if (tilepath->ExtraDataSize > 0) {
    tilepath->ExtraData = new char[tilepath->ExtraDataSize];
  }

  sz = fread(tilepath->Loq, sizeof(unsigned char), n_8, fp);
  if (sz!=n_8) { return -1; }

  sz = fread(tilepath->Span, sizeof(unsigned char), n_8, fp);
  if (sz!=n_8) { return -1; }

  sz = fread(tilepath->Cache, sizeof(uint64_t), n_32, fp);
  if (sz!=n_32) { return -1; }

  if (tilepath->Overflow) {
    sz = fread(tilepath->Overflow, sizeof(uint16_t), (size_t)(tilepath->NOverflow), fp);
    if (sz!=(size_t)(tilepath->NOverflow)) { return -1; }
  }

  if (tilepath->Overflow64) {
    sz = fread(tilepath->Overflow64, sizeof(uint64_t), (size_t)(tilepath->NOverflow64), fp);
    if (sz!=(size_t)(tilepath->NOverflow64)) { return -1; }
  }

  if (tilepath->ExtraData) {
    sz = fread(tilepath->ExtraData, sizeof(char), (size_t)(tilepath->ExtraDataSize), fp);
    if (sz!=(size_t)(tilepath->ExtraDataSize)) { return -1; }
  }

  //---

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileStepHomSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileVariantHomSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileNocSumHomSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileNocStartHomSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileNocLenHomSize = u64;

  //---

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileStepHetSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileVariantHetSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileNocSumHetSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileNocStartHetSize = u64;

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { return -1; }
  tilepath->LoqTileNocLenHetSize = u64;


  k = _read_enc_vector(tilepath->LoqTileStepHom, tilepath->LoqTileStepHomSize, fp);
  if (k<0) { return k; }
  k = _read_vlc_vector(tilepath->LoqTileVariantHom, tilepath->LoqTileVariantHomSize, fp);
  if (k<0) { return k; }
  k = _read_enc_vector(tilepath->LoqTileNocSumHom, tilepath->LoqTileNocSumHomSize, fp);
  if (k<0) { return k; }
  k = _read_vlc_vector(tilepath->LoqTileNocStartHom, tilepath->LoqTileNocStartHomSize, fp);
  if (k<0) { return k; }
  k = _read_vlc_vector(tilepath->LoqTileNocLenHom, tilepath->LoqTileNocLenHomSize, fp);
  if (k<0) { return k; }

  k = _read_enc_vector(tilepath->LoqTileStepHet, tilepath->LoqTileStepHetSize, fp);
  if (k<0) { return k; }
  k = _read_vlc_vector(tilepath->LoqTileVariantHet, tilepath->LoqTileVariantHetSize, fp);
  if (k<0) { return k; }
  k = _read_enc_vector(tilepath->LoqTileNocSumHet, tilepath->LoqTileNocSumHetSize, fp);
  if (k<0) { return k; }
  k = _read_vlc_vector(tilepath->LoqTileNocStartHet, tilepath->LoqTileNocStartHetSize, fp);
  if (k<0) { return k; }
  k = _read_vlc_vector(tilepath->LoqTileNocLenHet, tilepath->LoqTileNocLenHetSize, fp);
  if (k<0) { return k; }


}

cgf_t *cgft_read(FILE *fp) {
  int i, j, k, ch;
  char buf[1024];
  cgf_t *cgf=NULL;
  tilepath_t *tilepath;

  size_t sz;

  uint64_t u64;
  uint32_t u32;
  unsigned char ub[32];

  cgf = new cgf_t;

  sz = fread(&ub, sizeof(char), 8, fp);
  if (sz!=8) { goto cgft_read_error; }
  for (i=0; i<8; i++) {
    if (ub[i] != CGFT_MAGIC[i]) { goto cgft_read_error; }
    cgf->Magic[i] = ub[i];
  }

  sz = fread(&u32, sizeof(uint32_t), 1, fp);

  if (sz!=1) { goto cgft_read_error; }

  cgf->CGFVersion.clear();
  cgf->CGFVersion.reserve(u32);
  for (i=0; i<u32; i++) {
    ch = fgetc(fp);
    if (ch==EOF) { goto cgft_read_error; }
    cgf->CGFVersion += ch;
  }

  sz = fread(&u32, sizeof(uint32_t), 1, fp);
  if (sz!=1) { goto cgft_read_error; }

  cgf->LibraryVersion.clear();
  cgf->LibraryVersion.reserve(u32);
  for (i=0; i<u32; i++) {
    ch = fgetc(fp);
    if (ch==EOF) { goto cgft_read_error; }
    cgf->LibraryVersion += ch;
  }

  sz = fread(&u64, sizeof(uint64_t), 1, fp);
  if (sz!=1) { goto cgft_read_error; }
  cgf->PathCount = u64;

  sz = fread(&u32, sizeof(uint32_t), 1, fp);
  if (sz!=1) { goto cgft_read_error; }
  cgf->TileMap.clear();
  cgf->TileMap.reserve(u32);
  for (i=0; i<u32; i++) {
    ch = fgetc(fp);
    if (ch==EOF) { goto cgft_read_error; }
    cgf->TileMap += ch;
  }

  cgf->PathStructOffset.clear();
  for (i=0; i<(int)cgf->PathCount; i++) {
    sz = fread(&u64, sizeof(uint64_t), 1, fp);
    if (sz!=1) { goto cgft_read_error; }

    cgf->PathStructOffset.push_back(u64);
  }

  for (i=0; i<(int)cgf->PathCount; i++) {
    tilepath_t tp;
    cgf->Path.push_back(tp);
    k = cgft_tilepath_read(&(cgf->Path[i]), fp);
    if (k<0) { goto cgft_read_error; }
  }


  return cgf;

cgft_read_error:

  if (cgf) { delete cgf; }
  return NULL;
}

void bprint(FILE *fp, uint32_t x) {
  int i;

  for (i=31; i>=0; i--) {
    fprintf(fp, "%c", (x & (1<<i)) ? '1' : '.');
  }
}


void mk_vec_tilemap(std::vector< std::vector< std::vector<int> > > &vtm, const char *tm) {
  int i, j, k, ii, jj;
  char *chp;
  std::string s;

  std::vector< std::vector<int> > tm_entry;
  std::vector< std::vector<int> > x;
  std::vector< int > y;

  int entry_count=0;
  int tval;

  int enc_val=0;

  std::string s_tm = tm;

  std::stringstream line_stream(s_tm);
  std::vector<std::string> lines, alleles, variants, vals;
  std::string item, item1, item2, item0;

  while (std::getline(line_stream, item, '\n')) {
    lines.push_back(item);
    entry_count++;
    if (entry_count==16) { break; }
  }

  vtm.clear();
  for (i=0; i<lines.size(); i++) {
    std::stringstream ss(lines[i]);

    tm_entry.clear();

    alleles.clear();
    while (std::getline(ss, item, ':')) {
      alleles.push_back(item);
    }

    for (j=0; j<alleles.size(); j++) {

      y.clear();
      std::stringstream s0(alleles[j]);
      while (std::getline(s0, item0, ';')) {

        int eo=0;
        std::stringstream s1(item0);
        while (std::getline(s1, item1, '+')) {

          tval = (int)strtol(item1.c_str(), NULL, 16);
          if (eo==1) {
            for (ii=1; ii<tval; ii++) { y.push_back(-1); }
          } else {
            y.push_back(tval);
          }

          eo = 1-eo;
        }
      }
      tm_entry.push_back(y);
    }
    vtm.push_back(tm_entry);

    enc_val++;
  }

}

int cgft_output_band_format(cgf_t *cgf, tilepath_t *tilepath, FILE *fp) {
  int i, j, k, ii, jj;
  uint64_t u64;
  uint32_t u32;
  uint32_t canon_mask, loq_mask, hiq_mask, span_mask, xspan_mask, anchor_mask;
  uint32_t mask, cache_mask, lo_cache;
  uint32_t cache_ovf_mask;

  int n_8, n_32, n_8_q, n_8_r, n_32_q, n_32_r, n, n_ovf;
  int pos8, pos32;

  unsigned char *loq, *span;
  uint64_t *cache;
  uint16_t *ovf;
  uint32_t valid_mask;

  int loc_verbose = 0;

  std::vector< std::vector< std::vector<int> > > tilemap_vec;

  n = (int)tilepath->NTileStep;
  n_8 = ((n+7)/8);
  n_32 = ((n+31)/32);

  n_8_q = n/8;
  n_32_q = n/32;
  n_8_r = n%8;
  n_32_r = n%32;

  loq = tilepath->Loq;
  span = tilepath->Span;
  cache = tilepath->Cache;
  ovf = tilepath->Overflow;
  n_ovf = (int)tilepath->NOverflow;

  std::vector<int> variant_v[2];
  std::vector< std::vector<int> > noc_v[2];

  std::vector<int> v;

  int hexit[8], cur_hexit, n_hexit, n_cache_ovf;

  mk_vec_tilemap(tilemap_vec, cgf->TileMap.c_str());

  for (i=0; i<n; i++) {
    variant_v[0].push_back(-1);
    variant_v[1].push_back(-1);
    noc_v[0].push_back(v);
    noc_v[1].push_back(v);
  }

  int tilestep = 0;
  for (ii=0; ii<n_32_q; ii++) {
    loq_mask = loq[4*ii] | (loq[4*ii+1]<<8) | (loq[4*ii+2]<<16) | (loq[4*ii+3]<<24);
    span_mask = span[4*ii] | (span[4*ii+1]<<8) | (span[4*ii+2]<<16) | (span[4*ii+3]<<24);

    xspan_mask = ~span_mask;
    hiq_mask = ~loq_mask;
    cache_mask = (uint32_t)(cache[ii]>>32);
    canon_mask = cache_mask & xspan_mask & hiq_mask;
    anchor_mask = span_mask & hiq_mask & (~cache_mask);
    cache_ovf_mask = (anchor_mask & hiq_mask) | ((~span_mask) & (~canon_mask) & hiq_mask);

    lo_cache = (uint32_t)(cache[ii]&(0xffffffff));

    if (loc_verbose) {
      fprintf(fp, "[%8i (q:%i,r:%i)]:\n",
          ii*32, ii, 0);
      fprintf(fp, "    hiq: "); bprint(fp, hiq_mask); fprintf(fp, "\n");
      fprintf(fp, "   span: "); bprint(fp, span_mask); fprintf(fp, "\n");
      fprintf(fp, "  cache: "); bprint(fp, cache_mask); fprintf(fp, "\n");
      fprintf(fp, "  canon: "); bprint(fp, canon_mask); fprintf(fp, "\n");
      fprintf(fp, "  anchr: "); bprint(fp, anchor_mask); fprintf(fp, "\n");
      fprintf(fp, "   covf: "); bprint(fp, cache_ovf_mask); fprintf(fp, "\n");
    }

    for (i=0; i<8; i++) {
      hexit[i] = (int)((lo_cache & (0xf<<(4*i))) >> (4*i));
    }

    cur_hexit=0;
    n_hexit=0;
    n_cache_ovf = 0;
    for (i=0; i<32; i++) {
      if (cache_ovf_mask & (1<<i)) {

        if ((cur_hexit<8) && (hexit[cur_hexit] != 0xf)) {
          int hexit_val = hexit[cur_hexit];
          for (j=0; j<tilemap_vec[hexit_val][0].size(); j++) {
            int cur_tilestep = 32*ii + i + j;


            if (loc_verbose) {
              fprintf(fp, "  cur_tilestep %i : hexit_val %i -> %i %i\n",
                  cur_tilestep,
                  hexit_val,
                  tilemap_vec[hexit_val][0][j],
                  tilemap_vec[hexit_val][1][j] );
            }

            variant_v[0][cur_tilestep] = tilemap_vec[hexit_val][0][j];
            variant_v[1][cur_tilestep] = tilemap_vec[hexit_val][1][j];
          }
        }
        cur_hexit++;

        n_cache_ovf++;
        n_hexit++;
      }
    }
    if (n_hexit>8) { n_hexit=8; }


    for (i=0; i<32; i++) {

      if (canon_mask&(1<<i)) {
        variant_v[0][tilestep] = 0;
        variant_v[1][tilestep] = 0;

        if (loc_verbose) { fprintf(fp, "  tilestep %i -> 0 0\n", tilestep); }

      }
      else if (cache_ovf_mask & (1<<i)) {
        if (cur_hexit<n_hexit) {

          if (loc_verbose) {
            fprintf(fp, ">>> tilestep %i hexit %i\n", tilestep, hexit[cur_hexit]);
          }

        }
        cur_hexit++;
      }
      //else if (loq_mask & (1<<i)) { if (loc_verbose) { fprintf(fp, "  loq %i\n", tilestep); } }

      tilestep++;
    }

  }

  if (loc_verbose) {
    fprintf(fp, "...remainder\n");
  }

  if (n_32_r) {
    int tq_8 = (n_32_r+7)/8;

    loq_mask = 0;
    span_mask = 0;

    valid_mask = (0xffffffff>>(32-n_32_r));

    int base_pos = n_32_q * 32;
    for (i=0; i<n_32_r; i++) {
      int cur_pos = base_pos + i;

      int q = cur_pos / 8;
      int r = cur_pos % 8;

      if (loq[q] & (1<<r)) { loq_mask |= (1<<i); }
      if (span[q] & (1<<r)) { span_mask |= (1<<i); }

    }

    xspan_mask = (~span_mask) & valid_mask;
    hiq_mask = (~loq_mask) & valid_mask;
    cache_mask = (uint32_t)(cache[n_32_q]>>32);
    cache_mask &= valid_mask;
    canon_mask = cache_mask & xspan_mask & hiq_mask & valid_mask;
    anchor_mask = span_mask & hiq_mask & (~cache_mask) & valid_mask;
    cache_ovf_mask = (anchor_mask & hiq_mask) | ((~span_mask) & (~canon_mask) & hiq_mask);
    cache_ovf_mask &= valid_mask;

    lo_cache = (uint32_t)(cache[n_32_q]&(0xffffffff));

    if (loc_verbose) {
      fprintf(fp, "[%8i (q:%i,r:%i)]:\n",
          n_32_q*32, n_32_q, n_32_r);
      fprintf(fp, "  valid: "); bprint(fp, valid_mask); fprintf(fp, "\n");
      fprintf(fp, "    loq: "); bprint(fp, loq_mask); fprintf(fp, "\n");
      fprintf(fp, "    hiq: "); bprint(fp, hiq_mask); fprintf(fp, "\n");
      fprintf(fp, "   span: "); bprint(fp, span_mask); fprintf(fp, "\n");
      fprintf(fp, "  cache: "); bprint(fp, cache_mask); fprintf(fp, "\n");
      fprintf(fp, "  canon: "); bprint(fp, canon_mask); fprintf(fp, "\n");
      fprintf(fp, "  anchr: "); bprint(fp, anchor_mask); fprintf(fp, "\n");
      fprintf(fp, "   covf: "); bprint(fp, cache_ovf_mask); fprintf(fp, "\n");
    }

    for (i=0; i<8; i++) {
      hexit[i] = (int)((lo_cache & (0xf<<(4*i))) >> (4*i));
    }

    cur_hexit=0;
    n_hexit=0;
    n_cache_ovf = 0;
    for (i=0; i<n_32_r; i++) {
      if (cache_ovf_mask & (1<<i)) {

        if ((cur_hexit<8) && (hexit[cur_hexit] != 0xf)) {
          int hexit_val = hexit[cur_hexit];
          for (j=0; j<tilemap_vec[hexit_val][0].size(); j++) {
            int cur_tilestep = 32*n_32_q + i + j;

            if (loc_verbose) {
              fprintf(fp, "  cur_tilestep (r) %i : hexit_val %i -> %i %i\n",
                  cur_tilestep,
                  hexit_val,
                  tilemap_vec[hexit_val][0][j],
                  tilemap_vec[hexit_val][1][j] );
            }

            variant_v[0][cur_tilestep] = tilemap_vec[hexit_val][0][j];
            variant_v[1][cur_tilestep] = tilemap_vec[hexit_val][1][j];
          }
        }
        cur_hexit++;

        n_cache_ovf++;
        n_hexit++;
      }
    }
    if (n_hexit>8) { n_hexit=8; }


    for (i=0; i<n_32_r; i++) {

      if (canon_mask&(1<<i)) {
        variant_v[0][tilestep] = 0;
        variant_v[1][tilestep] = 0;

        if (loc_verbose) { fprintf(fp, "  tiletsep (r) %i -> 0 0\n", tilestep); }
      }
      else if (cache_ovf_mask & (1<<i)) {
        if (cur_hexit<n_hexit) {

          if (loc_verbose) {
            fprintf(fp, "  tilestep (r) %i hexit %i\n", tilestep, hexit[cur_hexit]);
          }
        }
        cur_hexit++;
      }

      tilestep++;
    }

  }

  if (loc_verbose) {
    fprintf(fp, "...\n");
  }

  // cache processed, now fill in with overflow

  for (i=0; i<n_ovf; i+=3) {

    tilestep = (int)ovf[i];
    int vara = (int)ovf[i+1];
    int varb = (int)ovf[i+2];

    if (vara >= OVF16_MAX) { vara = -1; }
    if (varb >= OVF16_MAX) { varb = -1; }

    variant_v[0][tilestep] = vara;
    variant_v[1][tilestep] = varb;
  }

  // finally, fill in with nocall

  int prev_noc_start = 0;
  for (i=0; i<tilepath->LoqTileStepHom.size(); i++) {
    tilestep = tilepath->LoqTileStepHom[i];
    int vara = tilepath->LoqTileVariantHom[2*i];
    int varb = tilepath->LoqTileVariantHom[2*i+1];

    if (vara==SPAN_SDSL_ENC_VAL) { vara = -1; }
    if (varb==SPAN_SDSL_ENC_VAL) { varb = -1; }

    variant_v[0][tilestep] = vara;
    variant_v[1][tilestep] = varb;

    if (loc_verbose) {
      fprintf(fp, "  loq hom %i -> %i %i\n", tilestep, vara, varb);
    }

    for (j=prev_noc_start; j<tilepath->LoqTileNocSumHom[i]; j++) {
      noc_v[0][tilestep].push_back( (int)tilepath->LoqTileNocStartHom[j] );
      noc_v[0][tilestep].push_back( (int)tilepath->LoqTileNocLenHom[j] );

      noc_v[1][tilestep].push_back( (int)tilepath->LoqTileNocStartHom[j] );
      noc_v[1][tilestep].push_back( (int)tilepath->LoqTileNocLenHom[j] );
    }
    prev_noc_start = (int)tilepath->LoqTileNocSumHom[i];
  }

  prev_noc_start=0;
  for (i=0; i<tilepath->LoqTileStepHet.size(); i++) {
    tilestep = tilepath->LoqTileStepHet[i];
    int vara = tilepath->LoqTileVariantHet[2*i];
    int varb = tilepath->LoqTileVariantHet[2*i+1];

    if (vara==SPAN_SDSL_ENC_VAL) { vara = -1; }
    if (varb==SPAN_SDSL_ENC_VAL) { varb = -1; }

    variant_v[0][tilestep] = vara;
    variant_v[1][tilestep] = varb;

    if (loc_verbose) {
      fprintf(fp, "  loq het %i -> %i %i\n", tilestep, vara, varb);
    }

    for (j=prev_noc_start; j<tilepath->LoqTileNocSumHet[2*i]; j++) {
      noc_v[0][tilestep].push_back( (int)tilepath->LoqTileNocStartHet[j] );
      noc_v[0][tilestep].push_back( (int)tilepath->LoqTileNocLenHet[j] );
    }
    prev_noc_start = (int)tilepath->LoqTileNocSumHet[2*i];

    for (j=prev_noc_start; j<tilepath->LoqTileNocSumHet[2*i+1]; j++) {
      noc_v[1][tilestep].push_back( (int)tilepath->LoqTileNocStartHet[j] );
      noc_v[1][tilestep].push_back( (int)tilepath->LoqTileNocLenHet[j] );
    }
    prev_noc_start = (int)tilepath->LoqTileNocSumHet[2*i+1];

  }

  // create out nocall vectors
  //

  for (i=0; i<2; i++) {
    fprintf(fp, "[");
    for (j=0; j<variant_v[i].size(); j++) {
      fprintf(fp, " %i", variant_v[i][j]);
    }
    fprintf(fp, "]\n");
  }

  for (i=0; i<2; i++) {
    fprintf(fp, "[");
    for (j=0; j<noc_v[i].size(); j++) {
      fprintf(fp, "[");
      for (k=0; k<noc_v[i][j].size(); k++) {
        fprintf(fp, " %i", noc_v[i][j][k]);
      }
      fprintf(fp, " ]");
    }
    fprintf(fp, "]\n");
  }

}

