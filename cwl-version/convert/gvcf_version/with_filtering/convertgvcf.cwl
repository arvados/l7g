$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    inputBinding:
      position: 1
  gffDir: 
    type: Directory
    inputBinding:
      position: 2
  gffPrefix:
    type: string
    inputBinding:
      position: 3
  ref:
    type: string
    inputBinding:
      position: 4
  reffa: 
    type: File
    inputBinding:
      position: 5
  afn:
    type: File
    inputBinding:
      position: 6
  aidx:
    type: File
    inputBinding:
      position: 7
  refM:
    type: string
    inputBinding:
      position: 8
  reffaM:
    type: File
    inputBinding:
      position: 9
  afnM:
    type: File
    inputBinding:
      position: 10
  aidxM:
    type: File
    inputBinding:
      position: 11
  seqidM:
    type: string 
    inputBinding:
      position: 12
  tagdir:
    type: File
    inputBinding:
      position: 13
  l7g:
    type: File
    inputBinding:
      position: 14
  pasta:
    type: File
    inputBinding:
      position: 15
  refstream:
    type: File
    inputBinding:
      position: 16
  tile_assembly:
    type: File
    inputBinding:
      position: 17
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "stage/*"
  out2:
    type: File[]
    outputBinding:
      glob: "indexed/*.gz*"

