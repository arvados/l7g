cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  script: File
  tileassembly_hg19: File
  tagset: File
  refFa_hg38: File

outputs:
  result:
    type: File[]
    outputSource: hg38_tileassembly/result

steps:
  hg38_tileassembly:
    run: build-hg38-tileassembly.cwl
    in:
      script: script
      tileassembly_hg19: tileassembly_hg19
      tagset: tagset
      refFa_hg38: refFa_hg38
    out: [result]
