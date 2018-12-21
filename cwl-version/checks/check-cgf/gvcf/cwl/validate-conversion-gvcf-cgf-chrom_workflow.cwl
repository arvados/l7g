cwlVersion: v1.0
class: Workflow
label: Workflow to Validate the the gVCF to cgf iconversion with the tile library
requirements:
  ScatterFeatureRequirement: {}
  InlineJavascriptRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  script:
    type: File
    label: Script that runs the workflow
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
    label: Chromosomes to analyze
  tileassembly:
    type: File
    label: Location of tile assembly file
  refFaFn:
    type: File
    label: Reference FASTA File
  gvcfPrefixes:
    type: string[]
    label: Prefixes of gVCFs
  gvcfSuffixes:
    type: string[]
    label: Suffixes of gVCFs

outputs:
  result:
    type: Directory
    outputSource: gather/out
    label: Directory of validated cgf files

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
