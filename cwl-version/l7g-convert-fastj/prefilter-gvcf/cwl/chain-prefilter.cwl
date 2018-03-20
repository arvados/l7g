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
  script_chain:
    type: File
    inputBinding:
      position: 1

  script_prefilter_consolidate:
    type: File
    inputBinding:
      position: 2

  script_prefilter:
    type: File
    inputBinding:
      position: 3

  script_filter_qual:
    type: File
    inputBinding:
      position: 4

  out_gvcf:
    type: string
    inputBinding:
      position: 5

  qual_cutoff:
    type: string
    inputBinding:
      position: 6

  in_gvcfs:
    type: File[]
    inputBinding:
      position: 7

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."


