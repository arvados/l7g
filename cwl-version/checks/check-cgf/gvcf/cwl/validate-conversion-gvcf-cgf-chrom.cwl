cwlVersion: v1.0
class: CommandLineTool
label: Validate the conversion of the gVCF to cgf
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
    label: Master script to run validation
    inputBinding:
      position: 1

  cgfDir:
    type: Directory
    label: Compact genome format directory
    inputBinding:
      position: 2
    secondaryFiles:
      - .cgf

  sglfDir:
    type: Directory
    label: Tile library directory
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
    label: Reference tile assembly file
    inputBinding:
      position: 6
    secondaryFiles:
      - .fwi
      - .gzi

  refFaFn:
    type: File
    label: Reference FASTA file
    inputBinding:
      position: 7
    secondaryFiles:
      - .fai
      - .gzi

  gvcfPrefix:
    type: string
    label: Prefix for gVCF files
    inputBinding:
      position: 8

  gvcfSuffix:
    type: string
    label: Suffix for gVCF files
    inputBinding:
      position: 9

  outfileName:
    type: string
    label: Name of output file
    inputBinding:
      position: 10

outputs:
  result:
    type: File
    label: Validation logs
    outputBinding:
      glob: "*output.log"
