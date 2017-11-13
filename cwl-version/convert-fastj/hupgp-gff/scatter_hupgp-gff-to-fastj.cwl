cwlVersion: v1.0
class: Workflow
requirements:
  ScatterFeatureRequirement: {}

inputs:
  script: File
  gffFns: File[]
  tagset: File
  tileassembly: File
  refFaFn: File

# Source 'result' of type {"items": {"type": "array", "items": "File"}, "type": "array"} is incompatible
# with sink 'outfiles' of type {"type": "array", "items": "File"}

outputs:
  outfiles:
    type: {"items": {"type": "array", "items": "File"}, "type": "array"}
    outputSource: convert/result

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

