cwlVersion: v1.1
class: Workflow
label: Convert gVCF to FASTA for gVCF with NON_REF
requirements:
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil

inputs:
  sampleid:
    type: string
    label: Sample ID
  vcf:
    type: File
    label: Input gVCF
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
  fixvcf-get_bed_varonlyvcf:
    run: fixvcf-get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      vcf: vcf
      gqcutoff: gqcutoff
      genomebed: genomebed
    out: [nocallbed, varonlyvcf]

  bcftools-consensus:
    run: bcftools-consensus.cwl
    scatter: haplotype
    in:
      sampleid: sampleid
      vcf: fixvcf-get_bed_varonlyvcf/varonlyvcf
      ref: ref
      haplotype: haplotypes
      mask: fixvcf-get_bed_varonlyvcf/nocallbed
    out: [fa]
