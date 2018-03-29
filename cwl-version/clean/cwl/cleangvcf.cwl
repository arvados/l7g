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
    coresMin: 1 
    coresMax: 1 
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    inputBinding:
      position: 1
  gvcfDir: 
    type: Directory
    inputBinding:
      position: 2
  gvcfPrefix:
    type: string
    inputBinding:
      position: 3
  cleanvcf:
    type: File
    inputBinding:
      position: 4 
outputs:
  out1:
    type: Directory 
    outputBinding:
      glob: "cleaned/*"
