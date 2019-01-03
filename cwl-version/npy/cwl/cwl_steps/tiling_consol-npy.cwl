$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    ramMin: 100000
    coresMin: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscriptmain:
    type: File
    default:
      class: File
      location: ../../src/allconsolCWL.sh
    inputBinding:
      position: 1
  indir:
    type: Directory
    inputBinding:
      position: 2
  outdir:
    type: string
    default: "outdir"
    inputBinding:
      position: 3
  outprefix:
    type: string
    default: "all"
    inputBinding:
      position: 4
  npyconsolfile:
    type: [File,string]
    default: "/usr/local/bin/npy-consolidate"
    inputBinding:
      position: 5
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: $(inputs.outdir)
