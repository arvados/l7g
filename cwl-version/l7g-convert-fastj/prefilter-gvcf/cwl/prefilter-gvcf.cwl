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
  gvcf:
    type: File
    inputBinding:
      position: 2
  ofn:
    type: string
    inputBinding:
      position: 3

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."


