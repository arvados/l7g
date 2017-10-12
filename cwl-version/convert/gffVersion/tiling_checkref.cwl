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
  inputdir: 
    type: Directory 
    inputBinding:
      position: 2
  ref:
    type: File
    inputBinding:
      position: 3 
  tag:
    type: string
    inputBinding:
      position: 4
  refM0:
    type: File
    inputBinding:
      position: 5
  tagM0:
    type: string
    inputBinding:
      position: 6
  refM1:
    type: File
    inputBinding:
      position: 7
  tagM1:
    type: string
    inputBinding:
      position: 8
  refM2:
    type: File
    inputBinding:
      position: 9
  tagM2:
    type: string
    inputBinding:
      position: 10
  whichever:
    type: File
    inputBinding:
      position: 11
  tempdir:
    type: string 
    inputBinding:      
      position: 12
  numberify:
    type: File
    inputBinding:
      position: 13
outputs:
  out1:
    type: Directory[]
    outputBinding:
      glob: "ChrMref*"
