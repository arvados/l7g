cwlVersion: v1.0
class: Workflow
label: Workflow to validate the the gVCF to cgf conversion
requirements:
  ScatterFeatureRequirement: {}
  InlineJavascriptRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  script:
    type: File
    label: Master script to run validation
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
  refFaFn:
    type: File
    label: Reference FASTA file
  gvcfPrefixes:
    type: string[]
    label: Prefix for gVCF files
  gvcfSuffixes:
    type: string[]
    label: Suffix for gVCF files

outputs:
  result:
    type: Directory
    outputSource: gather/out
    label: Directory of cgf validation logs

steps:
  cgf_gvcf_check:
    run: validate-conversion-gvcf-cgf-chrom.cwl
    scatter: [ chrom, gvcfPrefix, gvcfSuffix ]
    scatterMethod: dotproduct
    in:
      script: script
      cgfDir: cgfDir
      sglfDir: sglfDir
      gvcfDir: gvcfDir
      chrom: chroms
      refFaFn: refFaFn
      tileassembly: tileassembly
      gvcfPrefix: gvcfPrefixes
      gvcfSuffix: gvcfSuffixes
      outfileName:
        valueFrom: $(inputs.chrom + "-output.log")
    out: [result]
  gather:
    run: gather_validate-conversion-gvcf-cgf.cwl
    in:
      infiles: cgf_gvcf_check/result
    out: [out]
