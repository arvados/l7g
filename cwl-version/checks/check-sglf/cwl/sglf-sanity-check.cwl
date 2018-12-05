cwlVersion: v1.0
class: CommandLineTool
label: Check tile library (SGLFs) for correct formatting and spurious characters
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
    label: Bash master workflow directory
  sglfDir:
    type: Directory
    label: Tile Library Directory
    inputBinding:
      position: 2
  outFileName:
    type: string
    label: Name of output file, often includes chrom number
    inputBinding:
      position: 3

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."
