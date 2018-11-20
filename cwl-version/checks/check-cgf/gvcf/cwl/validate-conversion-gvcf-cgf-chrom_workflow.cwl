cwlVersion: v1.0
class: Workflow
label:
requirements:
  ScatterFeatureRequirement: {}
  InlineJavascriptRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  script: File
  cgfDir: Directory
  sglfDir: Directory
  gvcfDir: Directory
  chroms: string[]
  tileassembly: File
  refFaFn: File
  gvcfPrefixes: string[]
  gvcfSuffixes: string[]

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

