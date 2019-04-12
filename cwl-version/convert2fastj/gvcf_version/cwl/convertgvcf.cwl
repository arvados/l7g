$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
    ramMin: 13000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Master script to create a FastJ for each gVCF
    default:
      class: File
      location: ../src/convert2fastjCWL
    inputBinding:
      position: 1
  gffDir:
    type: Directory
    label: Directory of input gVCFs
    inputBinding:
      position: 2
  gffPrefix:
    type: string
    label: Prefix of all gVCF files
    inputBinding:
      position: 3
  ref:
    type: string
    label: Reference genome
    inputBinding:
      position: 4
  reffa:
    type: File
    label: Reference genome in FASTA format
    inputBinding:
      position: 5
  afn:
    type: File
    label: Compressed assembly fixed width file
    inputBinding:
      position: 6
  aidx:
    type: File
    label: Assembly index file
    inputBinding:
      position: 7
  refM:
    type: string
    label: Mitochondrial reference genome
    inputBinding:
      position: 8
  reffaM:
    type: File
    label: Reference mitochondrial genome in FASTA format
    inputBinding:
      position: 9
  afnM:
    type: File
    label: Compressed mitochondrial assembly fixed width file
    inputBinding:
      position: 10
  aidxM:
    type: File
    label: Mitochondrial assembly index file
    inputBinding:
      position: 11
  seqidM:
    type: string
    label: Mitochondrial naming scheme for gVCF
    inputBinding:
      position: 12
  tagdir:
    type: File
    label: Compressed tagset in FASTA format
    inputBinding:
      position: 13
  l7g:
    type: string
    label: Lightning application for parsing and searching assembly files
    default: "/usr/local/bin/l7g"
    inputBinding:
      position: 14
  pasta:
    type: string
    label: Tool for streaming and converting variant call formats
    default: "/usr/local/bin/pasta"
    inputBinding:
      position: 15
  refstream:
    type: string
    label: Tool for streaming and converting variant call formats
    default: "/usr/local/bin/refstream"
    inputBinding:
      position: 16
  tile_assembly:
    type: string
    label: Tool to extract information from the tile assembly files
    default: "/usr/local/bin/tile-assembly"
    inputBinding:
      position: 17
outputs:
  out1:
    type: Directory
    label: Directories of FastJs
    outputBinding:
      glob: "stage/*"
  out2:
    type: File[]
    label: Indexed and zipped gVCFs
    outputBinding:
      glob: "indexed/*.gz*"
