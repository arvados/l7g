$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Consolidates multiple NumPy arrays into one large array
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    ramMin: 100000
    coresMin: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Master script to consolidate tile path NumPy arrays
    default:
      class: File
      location: src/allconsolCWL.sh
    inputBinding:
      position: 1
  npydir:
    type: Directory
    label: Input directory
    inputBinding:
      position: 2
  outdir:
    type: string
    label: Name of output directory
    default: "outdir"
    inputBinding:
      position: 3
  outprefix:
    type: string
    label: Prefix for consolidated arrays
    default: "all"
    inputBinding:
      position: 4
  npyconsolfile:
    type: string
    label: Program to consolidated NumPy arrays
    default: "/usr/local/bin/npy-consolidate"
    inputBinding:
      position: 5
outputs:
  consolnpydir:
    type: Directory
    label: Consolidated NumPy arrays
    outputBinding:
      glob: $(inputs.outdir)
