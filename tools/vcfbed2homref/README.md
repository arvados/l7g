vcfbed2homref
===

A small utility to convert a VCF file and a BED file that represents homozygous reference
regions to VCF with homozygous reference regions called.

Quick Start
---

```
make
./vcfbed2homref -r testdata/ref.fa.gz -b testdata/small.bed testdata/small.vcf | egrep -v '^#' | head
1       238902  .       g       <NON-REF>       .       PASS    END=238915      GT      0/0
1       238919  .       a       <NON-REF>       .       PASS    END=238939      GT      0/0
1       239084  .       c       <NON-REF>       .       PASS    END=239113      GT      0/0
1       239230  .       c       <NON-REF>       .       PASS    END=239338      GT      0/0
1       239339  rs184451216     A       G       50      PASS    platforms=1;platformnames=10X;datasets=1;datasetnames=10XChromium;callsets=1;callsetnames=10XGATKhaplo;datasetsmissingcall=HiSeqPE300x,CGnormal,IonExome,SolidPE50x50bp,SolidSE75bp;callable=CS_10XGATKhaplo_callable;difficultregion=hg19_self_chain_split_withalts_gt10k      GT:DP:ADALL:AD:GQ:IGT:IPS:PS    0/1:0:0,0:0,0:99:0/1:.:.
1       239340  .       g       <NON-REF>       .       PASS    END=239431      GT      0/0
1       239434  .       a       <NON-REF>       .       PASS    END=239441      GT      0/0
1       239446  .       a       <NON-REF>       .       PASS    END=239481      GT      0/0
1       239482  rs201702841     G       T       50      PASS    platforms=2;platformnames=10X,Illumina;datasets=2;datasetnames=10XChromium,HiSeqPE300x;callsets=3;callsetnames=10XGATKhaplo,HiSeqPE300xGATK,HiSeqPE300xfreebayes;datasetsmissingcall=CGnormal,IonExome,SolidPE50x50bp,SolidSE75bp;callable=CS_10XGATKhaplo_callable;filt=CS_HiSeqPE300xGATK_filt,CS_HiSeqPE300xfreebayes_filt;difficultregion=hg19_self_chain_split_withalts_gt10k      GT:DP:ADALL:AD:GQ:IGT:IPS:PS    0/1:403:92,104:0,0:198:0/1:.:.
1       239624  .       t       <NON-REF>       .       PASS    END=239631      GT      0/0
```

This utility requires [htslib](https://github.com/samtools/htslib).

Usage
---

```
  vcfbed2homref [-h] [-v] [-V] [-N non-ref] [-b bedfile] [-r ref-fasta] [vcf_file] [out_vcf_file]

    [vcf_file]      VCF file (defaults stdin)
    [out_vcf_file]  output VCF file (defaults stdout)
    [-b bedfile]    bed file of homozygous ref sequences
    [-r ref-fasta]  reference FASTA file (indexed)
    [-s]            supress header (default output header)
    [-v]            verbose
    [-N non-ref]    "non ref" string to use (defaults to '<NON_REF>')
    [-V]            print version
    [-h]            Help (this screen)
```

Notes
---

The BED file is loaded into memory at runtime.
The VCF file is streamed in.
The reference file is accessed randomely by htslib.

* Reference FASTA file must be indexed (with `bgzip`)
* `testdata/small.vcf` and `testdata/small.bed` are small portions of the Genome in a Bottle `NA128178` that has both the [VCF](ftp://ftp-trace.ncbi.nlm.nih.gov/giab/ftp/release/NA12878_HG001/latest/GRCh37/HG001_GRCh37_GIAB_highconf_CG-IllFB-IllGATKHC-Ion-10X-SOLID_CHROM1-X_v.3.3.2_highconf_PGandRTGphasetransfer.vcf.gz) and [BED](ftp://ftp-trace.ncbi.nlm.nih.gov/giab/ftp/release/NA12878_HG001/latest/GRCh37/HG001_GRCh37_GIAB_highconf_CG-IllFB-IllGATKHC-Ion-10X-SOLID_CHROM1-X_v.3.3.2_highconf_nosomaticdel.bed) file available.
* The output VCF tries to leave the original lines untouched and only insert homozygous reference regions.

License
---

AGPLv3
