cwlVersion: v1.0
class: CommandLineTool
$namespaces:
  arv: "http://arvados.org/cwl#"
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: arv:RuntimeConstraints
    keep_cache: 10000

baseCommand: bash

inputs:

  script:
    type: File
    inputBinding:
      position: 1

  gvcfFn:
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

  refName:
    type: string
    inputBinding:
      position: 6

  outName:
    type: string
    inputBinding:
      position: 7

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."
