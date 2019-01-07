cwlVersion: v1.0
class: CommandLineTool
$namespaces:
  arv: "http://arvados.org/cwl#"
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
hints:
  arv:RuntimeConstraints:
    keep_cache: 10000
baseCommand: bash
inputs:
  script:
    type: File
    default:
      class: File
      location: ../src/sglf-sanity-check
    inputBinding:
      position: 1
  sglfDir:
    type: Directory
    inputBinding:
      position: 2
  outFileName:
    type: string
    inputBinding:
      position: 3
outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."
