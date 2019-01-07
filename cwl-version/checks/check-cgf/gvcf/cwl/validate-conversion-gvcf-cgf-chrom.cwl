cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
    ramMin: 65536
    ramMax: 131072
baseCommand: bash
inputs:
  script:
    type: File
    default:
      class: File
      location: ../src/verify-conversion-batch-gvcf-cgf_skip-empty-and-zero-tilepaths.sh
    inputBinding:
      position: 1
  cgfDir:
    type: Directory
    inputBinding:
      position: 2
    secondaryFiles:
      - .cgf
  sglfDir:
    type: Directory
    inputBinding:
      position: 3
    secondaryFiles:
      - .gz
  gvcfDir:
    type: Directory
    inputBinding:
      position: 4
    secondaryFiles:
      - .gz
      - .tbi
  chrom:
    type: string
    inputBinding:
      position: 5
  tileassembly:
    type: File
    inputBinding:
      position: 6
    secondaryFiles:
      - .fwi
      - .gzi
  refFaFn:
    type: File
    inputBinding:
      position: 7
    secondaryFiles:
      - .fai
      - .gzi
  gvcfPrefix:
    type: string
    inputBinding:
      position: 8
  gvcfSuffix:
    type: string
    inputBinding:
      position: 9
  outfileName:
    type: string
    inputBinding:
      position: 10
outputs:
  result:
    type: File
    outputBinding:
      glob: "*output.log"
