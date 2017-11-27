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
  gffFn:
    type: File
    inputBinding:
      position: 2
    secondaryFiles:
      - .tbi
      - .gzi
  tagset:
    type: File
    inputBinding:
      position: 3
    secondaryFiles:
      - .fai
      - .gzi
  tileassembly:
    type: File
    inputBinding:
      position: 4
    secondaryFiles:
      - .fwi
      - .gzi
  refFaFn:
    type: File
    inputBinding:
      position: 5
    secondaryFiles:
      - .fai
      - .gzi
  name:
    type: string
    default: ""
    inputBinding:
      position: 6

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."


