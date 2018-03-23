cwlVersion: v1.0
class: CommandLineTool
$namespaces:
  arv: "http://arvados.org/cwl#"
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: arv:RuntimeConstraints
    keep_cache: 10000

baseCommand: bash

inputs:
  script:
    type: File
    inputBinding:
      position: 1
  srcdir:
    type: Directory
    inputBinding:
      position: 2
  nppdir:
    type: Directory
    inputBinding:
      position: 3
  dstdir:
    type: string
    inputBinding:
      position: 4
  nthreads:
    type: string
    default: "."
    inputBinding:
      position: 5

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."

