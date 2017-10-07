cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement

baseCommand: bash

inputs:
  script:
    type: File
    inputBinding:
      position: 1
  tagset:
    type: File
    inputBinding:
      position: 2
  refFaFn:
    type: File
    inputBinding:
      position: 3
    secondaryFiles:
      - .fai
      - .gzi
  cytobandFn:
    type: File
    inputBinding:
      position: 4

outputs:
  result:
    type: File[]
    outputBinding:
      glob: "assembly.*"


