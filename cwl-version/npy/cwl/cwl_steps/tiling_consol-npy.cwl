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
    ramMin: 100000
    coresMin: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscriptmain:
    type: File?
    inputBinding:
      position: 1
    default:
      class: File
      location: ../../src/allconsolCWL.sh 
  indir:
    type: Directory 
    inputBinding:
      position: 2
  outdir:
    type: string? 
    inputBinding:
      position: 3
    default: "outdir"
  outprefix:
    type: string? 
    inputBinding:
      position: 4
    default: "all" 
  npyconsolfile:
    type: File?
    inputBinding:
      position: 5
    default:
      class: File
      location: ../../src/buildArvados/dest/npy-consolidate
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: $(inputs.outdir)
