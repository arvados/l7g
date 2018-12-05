cwlVersion: v1.0
class: Workflow
label:
requirements:
  ScatterFeatureRequirement: {}
  InlineJavascriptRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  script:
    type: File
    label: Bash script that runs the Workflow
  cgfDir:
    type: Directory
    label: Compact Genome Format Directory
  sglfDir:
    type: Directory
    label: Tile Library Directory
  gvcfDir:
    type: Directory
    label: gVCF Directory
  chroms:
    type: string[]
    label: Chromosomes to analyze
  tileassembly:
    type: File
  refFaFn:
    type: File
  gvcfPrefixes:
    type: string[]
    label: Prefixes to add to gVCF
  gvcfSuffixes:
    type: string[]
    label: Suffixes to add to gVCF

outputs:
  result:
    type: Directory
    outputSource: gather/out

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
