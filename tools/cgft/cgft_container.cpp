#include "cgft.hpp"

void cgft_create_container(FILE *fp, const char *tilemap) {
  uint32_t u32;
  uint64_t u64;
  int i, n;

  fwrite(CGFT_MAGIC, sizeof(char), 8, fp);

  u32 = (uint32_t)strlen(CGF_VERSION);
  fwrite(&u32, sizeof(uint32_t), 1, fp);
  fwrite(CGF_VERSION, sizeof(char), u32, fp);

  u32 = (uint32_t)strlen(CGLF_VERSION);
  fwrite(&u32, sizeof(uint32_t), 1, fp);
  fwrite(CGLF_VERSION, sizeof(char), u32, fp);

  u64 = 0;
  fwrite(&u64, sizeof(uint64_t), 1, fp);

  u32 = (uint32_t)strlen(tilemap);
  fwrite(&u32, sizeof(uint32_t), 1, fp);
  fwrite(tilemap, sizeof(char), u32, fp);

  // 0 PathStructOffset
  // no PathStruct...

}

int cgft_sanity(cgf_t *cgf) {
  int ii;
  int i, j, k;
  tilepath_t *tilepath;
  int n;

  for (ii=0; ii<cgf->Path.size(); ii++) {

    tilepath = &(cgf->Path[ii]);

    if ((2*tilepath->LoqTileStepHom.size()) != (tilepath->LoqTileVariantHom.size())) {
      return -1;
    }

    if (tilepath->LoqTileStepHom.size() != tilepath->LoqTileNocSumHom.size()) {
      return -2;
    }

    if (tilepath->LoqTileNocSumHom.size()>0) {
      n = (int)tilepath->LoqTileNocSumHom[ tilepath->LoqTileNocSumHom.size()-1 ];
      if (n != (int)tilepath->LoqTileNocStartHom.size()) {
        return -3;
      }

      if (n != (int)tilepath->LoqTileNocLenHom.size()) {
        return -4;
      }

    }

    if ((2*tilepath->LoqTileStepHet.size()) != (tilepath->LoqTileVariantHet.size())) {
      return -5;
    }

    if ((2*tilepath->LoqTileStepHet.size()) != (tilepath->LoqTileNocSumHet.size())) {
      return -6;
    }

    if (tilepath->LoqTileNocSumHet.size()>0) {
      n = (int)tilepath->LoqTileNocSumHet[ tilepath->LoqTileNocSumHet.size()-1 ];
      if (n != (int)tilepath->LoqTileNocStartHet.size()) {
        return -7;
      }

      if (n != (int)tilepath->LoqTileNocLenHet.size()) {
        return -8;
      }

    }

  }

  return 0;
}

void cgft_print_header(cgf_t *cgf) {
  int i, n;
  int tm_ele=0;

  printf("Magic:");
  for (i=0; i<8; i++) { printf(" %c", cgf->Magic[i]); }
  printf("\n");

  printf("CGFVersion: %s\n", cgf->CGFVersion.c_str());
  printf("LibraryVersion: %s\n", cgf->LibraryVersion.c_str());
  printf("PathCount: %llu\n", (long long unsigned int)cgf->PathCount);
  printf("TileMap(%i):\n", (int)cgf->TileMap.size());

  n = ( (cgf->TileMap.size() < 30) ? (int)(cgf->TileMap.size()) : 30 );
  n = (int)cgf->TileMap.size();

  printf("   [");
  for (i=0; i<n; i++) {
    if (cgf->TileMap[i] == '\n') {
      printf("]");
      tm_ele++;
      if ((tm_ele%16)==0) { printf("\n  "); }

      printf(" [");
    }
    else { printf("%c", cgf->TileMap[i]); }
  }
  printf("]\n");

  printf("PathStructOffset(%i):\n", (int)(cgf->PathStructOffset.size()));
  for (i=0; i<cgf->PathStructOffset.size(); i++) {
    if ((i>0) && (i%16)==0) { printf("\n"); }
    printf("  %llu", (long long unsigned int)cgf->PathStructOffset[i]);
  }
  printf("\n");

}

void cgft_print_tilepath(cgf_t *cgf, tilepath_t *tilepath) {
  int i, j, v;
  int n, n_bv, n_cache;
  uint64_t u64;
  uint32_t u32, t32;
  uint16_t u16;
  int prev_sum=0;
  int OVF_VAL =  0xffff;
  uint64_t OVF64_VAL =  0xffffffffffffffff;

  n_bv = (int)((tilepath->NTileStep+7)/8);
  n_cache = (int)((tilepath->NTileStep+31)/32);


  printf("TilePath: %i (%03x)\n", (int)tilepath->TilePath, (int)tilepath->TilePath);
  printf("Name: %s\n", tilepath->Name.c_str());
  printf("NTileStep: %i\n", (int)tilepath->NTileStep);
  printf("NOverflow: %i\n", (int)tilepath->NOverflow);
  printf("NOverflow64: %i\n", (int)tilepath->NOverflow64);
  printf("ExtraDataSize: %i\n", (int)tilepath->ExtraDataSize);

  printf("Cache(%i):\n", n_cache);
  printf("  ");
  for (i=0; i<n_cache; i++) {
    if ((i>0) && ((i%4)==0)) { printf("\n  "); }
    u64 = tilepath->Cache[i];
    u32 = (uint32_t)(u64>>32);
    printf(" %08x", u32);
    printf(" ");

    u32 = (uint32_t)(u64&0xffffffff);
    //for (j=0; j<(32/4); j++) {
    for (j=((32/4)-1); j>=0; j--) {
      t32 = (u32 & (0xf<<(4*j))) >> (4*j);
      printf("%x", t32);
    }
  }
  printf("\n");

  // 16 bit overflow
  //

  printf("Overflow(%i):\n", (int)tilepath->NOverflow);
  n = (int)(tilepath->NOverflow);
  printf(" ");
  for (i=0; i<n; i+=3) {
    if ((i>0) && ((i%16)==0)) { printf("\n "); }
    printf(" [%i]", (int)tilepath->Overflow[i]);
    printf("(");

    v = (int)tilepath->Overflow[i+1];
    if (v==OVF_VAL) { printf("."); }
    else { printf("%i", v); }

    printf(",");

    v = (int)tilepath->Overflow[i+2];
    if (v==OVF_VAL) { printf("."); }
    else { printf("%i", v); }

    printf(")");

  }
  printf("\n");

  // Overflow64
  //

  printf("Overflow64(%i):\n", (int)tilepath->NOverflow64);
  n = (int)(tilepath->NOverflow64);
  printf(" ");
  for (i=0; i<n; i+=3) {
    if ((i>0) && ((i%16)==0)) { printf("\n "); }
    printf(" [%i]", (int)tilepath->Overflow64[i]);
    printf("(");

    if (tilepath->Overflow64[i+1]==OVF64_VAL) { printf("."); }
    else { printf("%llu", (long long unsigned)tilepath->Overflow64[i+1]); }

    printf(",");

    if (tilepath->Overflow64[i+2]==OVF64_VAL) { printf("."); }
    else { printf("%llu", (long long unsigned)tilepath->Overflow64[i+2]); }


    printf(")");

  }
  printf("\n");

  // Low quality
  //

  printf("Loq(%i):\n", n_bv);
  printf("  ");
  for (i=0; i<n_bv; i++) {
    if ((i>0) && ((i%32)==0)) { printf("\n  "); }
    printf(" %02x", (int)(tilepath->Loq[i]));
  }
  printf("\n");

  printf("Span(%i):\n", n_bv);
  printf("  ");
  for (i=0; i<n_bv; i++) {
    if ((i>0) && ((i%32)==0)) { printf("\n  "); }
    printf(" %02x", (int)(tilepath->Span[i]));
  }
  printf("\n");

  // Extra Data
  //

  printf("ExtraData(%i):\n", (int)tilepath->ExtraDataSize);
  n = (int)tilepath->ExtraDataSize;
  for (i=0; i<n; i++) {
    if ((i>0) && ((i%32)==0)) { printf("\n  "); }
    printf(" %2x", (int)(tilepath->ExtraData[i]));
  }
  printf("\n");

  printf("HomSizes: Step:%i, Var:%i, NocSum:%i, NocStart:%i, NocLen:%i\n",
      (int)tilepath->LoqTileStepHomSize,
      (int)tilepath->LoqTileVariantHomSize,
      (int)tilepath->LoqTileNocSumHomSize,
      (int)tilepath->LoqTileNocStartHomSize,
      (int)tilepath->LoqTileNocLenHomSize);

  printf("HetSizes: Step:%i, Var:%i, NocSum:%i, NocStart:%i, NocLen:%i\n",
      (int)tilepath->LoqTileStepHetSize,
      (int)tilepath->LoqTileVariantHetSize,
      (int)tilepath->LoqTileNocSumHetSize,
      (int)tilepath->LoqTileNocStartHetSize,
      (int)tilepath->LoqTileNocLenHetSize);

  // Hom loq quaility
  //

  printf("LoqHom(%i,%i,%i,%i,%i):\n",
      (int)tilepath->LoqTileStepHom.size(),
      (int)tilepath->LoqTileVariantHom.size(),
      (int)tilepath->LoqTileNocSumHom.size(),
      (int)tilepath->LoqTileNocStartHom.size(),
      (int)tilepath->LoqTileNocLenHom.size() );

  n = (int)tilepath->LoqTileStepHom.size();
  printf(" ");
  prev_sum = 0;
  for (i=0; i<n; i++) {
    if ((i>0) && ((i%4)==0)) { printf("\n "); }
    printf(" {");
    printf("%i(%i,%i)[%i]:",
        (int)tilepath->LoqTileStepHom[i],
        (int)tilepath->LoqTileVariantHom[2*i],
        (int)tilepath->LoqTileVariantHom[2*i+1],
        (int)tilepath->LoqTileNocSumHom[i]);
    for (j=prev_sum; j<tilepath->LoqTileNocSumHom[i]; j++) {
      if (j>prev_sum) { printf(" "); }
      printf("[%i+%i]",
          (int)tilepath->LoqTileNocStartHom[j],
          (int)tilepath->LoqTileNocLenHom[j]);
    }
    prev_sum = (int)tilepath->LoqTileNocSumHom[i];
    printf("}");
  }
  printf("\n");

  // Het loq quaility
  //

  printf("LoqHet(%i,%i,%i,%i,%i):\n",
      (int)tilepath->LoqTileStepHet.size(),
      (int)tilepath->LoqTileVariantHet.size(),
      (int)tilepath->LoqTileNocSumHet.size(),
      (int)tilepath->LoqTileNocStartHet.size(),
      (int)tilepath->LoqTileNocLenHet.size() );

  n = (int)tilepath->LoqTileStepHet.size();
  printf(" ");
  prev_sum = 0;
  for (i=0; i<n; i++) {
    if ((i>0) && ((i%4)==0)) { printf("\n "); }
    printf(" {");
    printf("%i(%i,%i)[%i,%i]:",
        (int)tilepath->LoqTileStepHet[i],
        (int)tilepath->LoqTileVariantHet[2*i],
        (int)tilepath->LoqTileVariantHet[2*i+1],
        (int)tilepath->LoqTileNocSumHet[2*i],
        (int)tilepath->LoqTileNocSumHet[2*i+1]);

    for (j=prev_sum; j<tilepath->LoqTileNocSumHet[2*i]; j++) {
      if (j>prev_sum) { printf(" "); }
      printf("[%i+%i]",
          (int)tilepath->LoqTileNocStartHet[j],
          (int)tilepath->LoqTileNocLenHet[j]);
    }
    prev_sum = (int)tilepath->LoqTileNocSumHet[2*i];

    for (j=prev_sum; j<tilepath->LoqTileNocSumHet[2*i+1]; j++) {
      if (j>prev_sum) { printf(" "); }

      if (j >= (int)tilepath->LoqTileNocStartHet.size()) {
        printf("SANITY: j (%i) >= LoqTileNocStartHet (%i)\n", j, (int)tilepath->LoqTileNocStartHet.size());
      }

      if (j >= (int)tilepath->LoqTileNocLenHet.size()) {
        printf("SANITY: j (%i) >= LoqTileNocLenHet (%i)\n", j, (int)tilepath->LoqTileNocLenHet.size());
      }

      printf("[%i+%i]",
          (int)tilepath->LoqTileNocStartHet[j],
          (int)tilepath->LoqTileNocLenHet[j]);
    }
    prev_sum = (int)tilepath->LoqTileNocSumHet[2*i+1];

    printf("}");
  }
  printf("\n");


  printf("\n");

}
