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
    label: Script that runs the workflow
    inputBinding:
      position: 1

  cgfDir:
    type: Directory
    label: Compact genome format (cgf) directory
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
    label: Tool to extract information from the tile assembly files
    inputBinding:
      position: 6
    secondaryFiles:
      - .fwi
      - .gzi

  refFaFn:
    type: File
    label: Reference FASTA File
    inputBinding:
      position: 7
    secondaryFiles:
      - .fai
      - .gzi

  gvcfPrefix:
    type: string
    label: Prefixes of gVCFs
    inputBinding:
      position: 8

  gvcfSuffix:
    type: string
    label: Suffixes of gVCFs
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
