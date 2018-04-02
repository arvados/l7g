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
  tagset: File
  refFaFn : File
  cytobandFn: File

outputs:
  result:
    type: File[]
    outputSource: tile_liftover/result

steps:
  tile_liftover:
    run: tile-liftover.cwl
    in:
      script: script
      tagset: tagset
      refFaFn : refFaFn
      cytobandFn: cytobandFn
    out: [result]
