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
  tilepath:
    type: string
    inputBinding:
      position: 2
  fastj_base_dir:
    type: Directory
    inputBinding:
      position: 3
  tagset:
    type: File
    inputBinding:
      position: 4
    secondaryFiles:
      - .fai
      - .gzi
  outdir:
    type: string
    default: "."
    inputBinding:
      position: 5

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."
