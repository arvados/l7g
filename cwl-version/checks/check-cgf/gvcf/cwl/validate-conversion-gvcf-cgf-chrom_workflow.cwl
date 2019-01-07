cwlVersion: v1.0
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  cgfDir: Directory
  sglfDir: Directory
  gvcfDir: Directory
  chroms: string[]
  tileassembly: File
  refFaFn: File
  gvcfPrefix: string
  gvcfSuffixes: string[]

outputs:
  result:
    type: Directory
    outputSource: gather/out

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
