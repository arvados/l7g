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
    outputSource: g1k_v37_tile_assembly/result

steps:
  g1k_v37_tile_assembly:
    run: build-human_g1k_v37-tile-assembly.cwl
    in:
      script: script
      tagset: tagset
      refFaFn : refFaFn
      cytobandFn: cytobandFn
    out: [result]
