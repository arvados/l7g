cwlVersion: v1.0
class: Workflow
label: Workflow to validate the the gVCF to cgf conversion
requirements:
  ScatterFeatureRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  cgfDir:
    type: Directory
    label: Compact genome format directory
  sglfDir:
    type: Directory
    label: Tile library directory
  gvcfDir:
    type: Directory
    label: gVCF directory
  chroms:
    type: string[]
    label: Chromosomes to analyze
  tileassembly:
    type: File
    label: Reference tile assembly file
    secondaryFiles:
      - .fwi
      - .gzi
  refFaFn:
    type: File
    label: Reference FASTA file
    secondaryFiles:
      - .fai
      - .gzi
  gvcfPrefix:
    type: string
    label: Prefix for gVCF files
  gvcfSuffixes:
    type: string[]
    label: Suffixes for gVCF files

outputs:
  result:
    type: Directory
    outputSource: gather/out
    label: Directory of cgf validation logs

steps:
  cgf_gvcf_check:
    run: validate-conversion-gvcf-cgf-chrom.cwl
    scatter: [chrom,gvcfSuffix]
    scatterMethod: dotproduct
    in:
      cgfDir: cgfDir
      sglfDir: sglfDir
      gvcfDir: gvcfDir
      chrom: chroms
      refFaFn: refFaFn
      tileassembly: tileassembly
      gvcfPrefix: gvcfPrefix
      gvcfSuffix: gvcfSuffixes
      outfileName:
        valueFrom: $(inputs.chrom)-output.log
    out: [result]
  gather:
    run: gather_validate-conversion-gvcf-cgf.cwl
    in:
      infiles: cgf_gvcf_check/result
    out: [out]
