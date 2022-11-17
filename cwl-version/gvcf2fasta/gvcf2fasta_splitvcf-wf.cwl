cwlVersion: v1.1
class: Workflow
label: Convert gVCF to FASTA for gVCF split by chromosome
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
  splitvcfdir:
    type: Directory
    label: Input directory of split gVCFs
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering
  genomebed:
    type: File
    label: Whole genome BED
  ref:
    type: File
    label: Reference FASTA

outputs:
  fas:
    type: File[]
    label: Output pair of FASTAs
    outputSource: bcftools-consensus/fas

steps:
  concat-get_bed_varonlyvcf:
    run: concat-get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      splitvcfdir: splitvcfdir
      gqcutoff: gqcutoff
      genomebed: genomebed
    out: [nocallbed, varonlyvcf]

  bcftools-consensus:
    run: bcftools-consensus.cwl
    in:
      sampleid: sampleid
      vcf: concat-get_bed_varonlyvcf/varonlyvcf
      ref: ref
      mask: concat-get_bed_varonlyvcf/nocallbed
    out: [fas]
