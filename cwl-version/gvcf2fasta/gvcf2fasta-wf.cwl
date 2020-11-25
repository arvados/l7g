cwlVersion: v1.1
class: Workflow
label: Convert gVCF to FASTA
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
  vcf:
    type: File
    label: Input gVCF
    secondaryFiles: [.tbi]
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
  get_bed_varonlyvcf:
    run: get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      vcf: vcf
      genomebed: genomebed
    out: [nocallbed, varonlyvcf]

  bcftools-consensus:
    run: bcftools-consensus.cwl
    scatter: haplotype
    in:
      sampleid: sampleid
      vcf: get_bed_varonlyvcf/varonlyvcf
      ref: ref
      haplotype: haplotypes
      mask: get_bed_varonlyvcf/nocallbed
    out: [fa]
