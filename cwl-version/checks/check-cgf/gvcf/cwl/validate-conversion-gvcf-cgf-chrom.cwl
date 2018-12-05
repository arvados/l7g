cwlVersion: v1.0
class: CommandLineTool
label: Validate the conversion of the gVCF to cgf against the SGLF (Tile Library)
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
    ramMin: 65536
    ramMax: 131072
  - class: InlineJavascriptRequirement

baseCommand: bash

inputs:

  script:
    type: File
    inputBinding:
      position: 1
    label: Bash script that runs the Workflow

  cgfDir:
    type: Directory
    label: Compact Genome Format Directory
    inputBinding:
      position: 2
    secondaryFiles:
      - .cgf

  sglfDir:
    type: Directory
    label: Tile Library Directory
    inputBinding:
      position: 3
    secondaryFiles:
      - .gz

  gvcfDir:
    type: Directory
    label: gVCF Directory
    inputBinding:
      position: 4
    secondaryFiles:
      - .gz
      - .tbi

  chrom:
    type: string
    label: Chromosomes to analyze
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
    label: Prefixes to add to gVCF
    inputBinding:
      position: 8

  gvcfSuffix:
    type: string
    label: Suffixes to add to gVCF
    inputBinding:
      position: 9

  outfileName:
    type: string
    label: Name of output file, often includes chrom number
    inputBinding:
      position: 10


outputs:

  result:
    type: File
    outputBinding:
      glob: "*output.log"
