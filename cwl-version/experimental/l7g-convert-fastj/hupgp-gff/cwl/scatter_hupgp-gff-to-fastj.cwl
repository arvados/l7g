cwlVersion: v1.0
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  InlineJavascriptRequirement: {}

inputs:
  script: File
  gffFns: File[]
  tagset: File
  tileassembly: File
  refFaFn: File

outputs:
  outfiles:
    type: Directory
    outputSource: gather/out

steps:
  convert:
    run: hupgp-gff-to-fastj.cwl
    scatter: gffFn
    in:
      script: script
      gffFn: gffFns
      tagset: tagset
      tileassembly: tileassembly
      refFaFn: refFaFn
    out: [result]
  gather:
    run: gather_hupgp-gff-to-fastj.cwl
    in:
      indirs: convert/result
    out: [out]

