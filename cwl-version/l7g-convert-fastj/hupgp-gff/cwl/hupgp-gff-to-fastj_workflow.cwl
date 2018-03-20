cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

baseCommand: bash

inputs:
  script: File
  gffFn: File
  tagset: File
  tileassembly: File
  refFaFn: File

outputs:
  result:
    type: File[]
    outputSource: hupgp_gff_to_fastj/result

steps:
  hupgp_gff_to_fastj:
    run: hupgp-gff-to-fastj.cwl
    in:
      script: script
      gffFn: gffFn
      tagset: tagset
      tileassembly: tileassembly
      refFaFn: refFaFn
    out: [result]
