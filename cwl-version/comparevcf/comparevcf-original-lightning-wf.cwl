cwlVersion: v1.2
class: Workflow
requirements:
  StepInputExpressionRequirement: {}
inputs:
  sampleid: string
  vcf: File
  nocallbed: File
  lightningvcf: File
  sdf: Directory
  bashscript: File
  header: File

outputs:
  modifiedoriginalvcf:
    type: File
    outputSource: change-GT_original/modifiedvcf
  modifiedlightningvcf:
    type: File
    outputSource: change-GT_lightning/modifiedvcf
  evaldir:
    type: Directory
    outputSource: rtg-vcfeval/evaldir

steps:
  rtg-vcffilter:
    run: ../imputation/rtg-vcffilter.cwl
    in:
      sample: sampleid
      vcf: vcf
      excludebed: nocallbed
    out: [filteredvcf]

  change-GT_original:
    run: change-GT.cwl
    in:
      sampleid: sampleid
      suffix:
        valueFrom: "original"
      vcf: rtg-vcffilter/filteredvcf
      header: header
      bashscript: bashscript
    out: [modifiedvcf]

  change-GT_lightning:
    run: change-GT.cwl
    in:
      sampleid: sampleid
      suffix:
        valueFrom: "lightning"
      vcf: lightningvcf
      header: header
      bashscript: bashscript
    out: [modifiedvcf]

  rtg-vcfeval:
    run: rtg-vcfeval.cwl
    in:
      baselinevcf: change-GT_original/modifiedvcf
      callsvcf: change-GT_lightning/modifiedvcf
      sdf: sdf
    out: [evaldir]
