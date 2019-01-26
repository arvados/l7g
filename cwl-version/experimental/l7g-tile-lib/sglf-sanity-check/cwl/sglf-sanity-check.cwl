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
    label: Bash script that checks for span and formatting errors in tile library
  sglfDir:
    type: Directory
    inputBinding:
      position: 2
    label: The directory conating tile library
  outFileName:
    type: string
    inputBinding:
      position: 3
    label: Name of output of the tile library sanity formatting checks

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."
