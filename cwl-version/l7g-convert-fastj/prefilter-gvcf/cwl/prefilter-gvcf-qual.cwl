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
  in_gvcf:
    type: File
    inputBinding:
      position: 2
  qual_cutoff:
    type: string
    inputBinding:
      position: 3
  out_gvcf:
    type: string
    inputBinding:
      position: 4

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."


