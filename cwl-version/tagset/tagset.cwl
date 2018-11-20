cwlVersion: v1.0
class: CommandLineTool
label: Output a FASTA file for the tagset
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
  cytobandFn:
    type: File
    inputBinding:
      position: 2
  bigwigFn:
    type: File
    inputBinding:
      position: 3
  refFaFn:
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
      glob: "*.fa.gz*"


