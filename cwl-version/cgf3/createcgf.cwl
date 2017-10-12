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
    ramMin: 14336 
    ramMax: 14336 
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    inputBinding:
      position: 1
  fjdir:
    type: Directory 
    inputBinding:
      position: 2
  cgft: 
    type: File 
    inputBinding:
      position: 3 
  fjt:
    type: File
    inputBinding:
      position: 4
  cglf:
    type: Directory
    inputBinding:
      position: 5 
outputs:    
  out1: 
    type: File 
    outputBinding:
      glob: "data/*.cgf"

