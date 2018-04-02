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
  tileassembly_hg19:
    type: File
    inputBinding:
      position: 2
    secondaryFiles:
      - .gzi
      - .fwi
  tagset:
    type: File
    inputBinding:
      position: 3
    secondaryFiles:
      - .fai
      - .gzi
  refFa_hg38:
    type: File
    inputBinding:
      position: 4
    secondaryFiles:
      - .fai
      - .gzi

outputs:
  result:
    type: File[]
    outputBinding:
      glob: "assembly.*"


