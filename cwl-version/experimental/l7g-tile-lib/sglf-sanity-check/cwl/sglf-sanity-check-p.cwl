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
  idir:
    type: Directory
    inputBinding:
      position: 2
  outFileName:
    type: string
    inputBinding:
      position: 3
  ncore:
    type: string
    inputBinding:
      position: 4

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."

