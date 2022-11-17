cwlVersion: v1.1
class: Workflow
label: Convert gVCF to FASTA for gVCF tar split by chromosome
requirements:
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000

inputs:
  sampleid:
    type: string
    label: Sample ID
  vcftar:
    type: File
    label: Input gVCF tar
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering
  genomebed:
    type: File
    label: Whole genome BED
  ref:
    type: File
    label: Reference FASTA
  haplotypes:
    type: int[]
    label: Haplotypes of sample
    default: [1, 2]

outputs:
  fas:
    type: File[]
    label: Output pair of FASTAs
    outputSource: bcftools-consensus/fa

steps:
  untar-concat-get_bed_varonlyvcf:
    run: untar-concat-get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      vcftar: vcftar
      gqcutoff: gqcutoff
      genomebed: genomebed
    out: [nocallbed, varonlyvcf]

  bcftools-consensus:
    run: bcftools-consensus.cwl
    scatter: haplotype
    in:
      sampleid: sampleid
      vcf: untar-concat-get_bed_varonlyvcf/varonlyvcf
      ref: ref
      haplotype: haplotypes
      mask: untar-concat-get_bed_varonlyvcf/nocallbed
    out: [fa]
