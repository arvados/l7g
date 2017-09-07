#include "cgft.hpp"

void cgft_tilepath_init(tilepath_t &tpath, uint64_t tilepath_name) {
  int i, j, k, n;
  char buf[1024];

  tpath.TilePath = tilepath_name;
  sprintf(buf, "%03x", (int)tilepath_name);
  tpath.NTileStep = 0;
  tpath.NOverflow = 0;
  tpath.NOverflow64 = 0;
  tpath.ExtraDataSize = 0;

  tpath.Loq = NULL;
  tpath.Span = NULL;
  tpath.Cache = NULL;

  tpath.Overflow = NULL;
  tpath.Overflow64 = NULL;

  tpath.ExtraData = NULL;

  /*
  tpath.LoqTileStepHom.resize(0);
  tpath.LoqTileVariantHom.clear();
  tpath.LoqTileNocSumHom.clear();
  tpath.LoqTileNocNocStartHom.clear();
  tpath.LoqTileNocNocLenHom.clear();

  tpath.LoqTileStepHet.clear();
  tpath.LoqTileVariantHet.clear();
  tpath.LoqTileNocSumHet.clear();
  tpath.LoqTileNocNocStartHet.clear();
  tpath.LoqTileNocNocLenHet.clear();
  */
}
