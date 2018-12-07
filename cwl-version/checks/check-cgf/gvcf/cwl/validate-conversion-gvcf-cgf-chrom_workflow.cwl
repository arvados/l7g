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
    label: Script that runs the Workflow
  cgfDir:
    type: Directory
    label: Compact genome format (cgf) directory
  sglfDir:
    type: Directory
    label: Tile library directory
  gvcfDir:
    type: Directory
    label: gVCF directory
  chroms:
    type: string[]
    label: Arrray of chromosomes to analyze
  tileassembly:
    type: File
    label: Tool to extract information from the tile assembly files
  refFaFn:
    type: File
    label: Reference Fasta File
  gvcfPrefixes:
    type: string[]
    label: Arrray of gVCF prefixes
  gvcfSuffixes:
    type: string[]
    label: Arrray of gVCF suffixes

outputs:
  result:
    type: Directory
    outputSource: gather/out
    label: Directory of validated cgfs

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
