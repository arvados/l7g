cwlVersion: v1.1
class: Workflow
requirements:
  StepInputExpressionRequirement: {}
inputs:
  sampleid:
    type: string
  vcf:
    type: File
  gqcutoff:
    type: int
  genomebed:
    type: File
  lightningvcfdir:
    type: Directory
  header:
    type: File
  sdf:
    type: Directory

outputs:
  modifiedoriginalvcf:
    type: File
    outputSource: change-header-GT_original/modifiedvcf
  modifiedlightningvcf:
    type: File
    outputSource: change-header-GT_lightning/modifiedvcf
  evaldir:
    type: Directory
    outputSource: rtg-vcfeval/evaldir

steps:
  fixvcf-get_bed_varonlyvcf:
    run: ../gvcf2fasta/fixvcf-get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      vcf: vcf
      gqcutoff: gqcutoff
      genomebed: genomebed
    out: [nocallbed, varonlyvcf]

  bedtools-intersect:
    run: bedtools-intersect.cwl
    in:
      sampleid: sampleid
      vcf: fixvcf-get_bed_varonlyvcf/varonlyvcf
      genomebed: genomebed
      nocallbed: fixvcf-get_bed_varonlyvcf/nocallbed
    out: [intersectvcf]

  concatenate:
    run: concatenate.cwl
    in:
      sampleid: sampleid
      vcfdir: lightningvcfdir
    out: [vcf]

  change-header-GT_original:
    run: change-header-GT.cwl
    in:
      sampleid: sampleid
      suffix:
        valueFrom: "original"
      header: header
      vcf: bedtools-intersect/intersectvcf
    out: [modifiedvcf]

  change-header-GT_lightning:
    run: change-header-GT.cwl
    in:
      sampleid: sampleid
      suffix:
        valueFrom: "lightning"
      header: header
      vcf: concatenate/vcf
    out: [modifiedvcf]

  rtg-vcfeval:
    run: rtg-vcfeval.cwl
    in:
      baselinevcf: change-header-GT_original/modifiedvcf
      callsvcf: change-header-GT_lightning/modifiedvcf
      sdf: sdf
    out: [evaldir]
