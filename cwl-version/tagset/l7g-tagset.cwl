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
  cytobandFn: File
  bigwigFn: File
  refFaFn : File
  script: File

outputs:
  result:
    type: File[]
    outputSource: tagset/result

steps:
  tagset:
    run: tagset.cwl
    in:
      cytobandFn: cytobandFn
      bigwigFn: bigwigFn
      refFaFn: refFaFn
      script: script
    out: [result]

