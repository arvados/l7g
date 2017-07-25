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
  bashscript:
    type: File
    inputBinding:
      position: 2
  fastj2cgflib:
    type: File 
    inputBinding:
      position: 3
  datadir: 
    type: Directory 
    inputBinding:
      position: 4 
  verbose_tagset:
    type: File
    inputBinding:
      position: 5 
  tagset:
    type: File
    inputBinding:
      position: 6 
  numthreads:
    type: string
    inputBinding:
      position: 7 
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "*lib*"

