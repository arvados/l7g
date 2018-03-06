cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement

baseCommand: python

inputs:
  script:
    type: File
    inputBinding:
      position: 1
  tileassembly:
    type: File
    inputBinding:
      position: 2
    secondaryFiles:
      - .gzi
      - .fwi
  out_fn:
    type: string
    inputBinding:
      position: 3

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."

