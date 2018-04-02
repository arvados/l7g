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
    ramMin: 130000
    coresMin: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscriptmain:
    type: File
    inputBinding:
      position: 1
  indir:
    type: Directory 
    inputBinding:
      position: 2
  outdir:
    type: string 
    inputBinding:
      position: 3
  outprefix:
    type: string 
    inputBinding:
      position: 4 
  npyconsolfile:
    type: File 
    inputBinding:
      position: 5
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: $(inputs.outdir)
