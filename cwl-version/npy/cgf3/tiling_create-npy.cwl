$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: pythontools 
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    coresMin: 16
    coresMax: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscriptmain:
    type: File
    inputBinding:
      position: 1
  cgft:
    type: File 
    inputBinding:
      position: 2
  cgfdirectory:
    type: Directory
    inputBinding:
      position: 3
  band2matrix:
    type: File
    inputBinding:
      position: 4 
  cnvrt2hiq:
    type: File
    inputBinding:
      position: 5
  makelist:
    type: File
    inputBinding:
      position: 6
  nthreads:
    type: string
    inputBinding:
      position: 7
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "npy"
  out2:
    type: Directory
    outputBinding:
      glob: "npy-hiq"
