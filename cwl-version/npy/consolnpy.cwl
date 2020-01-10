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
    coresMin: 16
    ramMin: 180000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  bashscript:
    type: File
    label: Master script to consolidate tile path NumPy arrays
    default:
      class: File
      location: src/allconsolCWL.sh
  npydir:
    type: Directory
    label: Input directory
  outdir:
    type: string
    label: Name of output directory
    default: "outdir"
  outprefix:
    type: string
    label: Prefix for consolidated arrays
    default: "all"
  npyconsolfile:
    type: string
    label: Program to consolidated NumPy arrays
    default: "/usr/local/bin/npy-consolidate"
outputs:
  consolnpydir:
    type: Directory
    label: Consolidated NumPy arrays
    outputBinding:
      glob: $(inputs.outdir)
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.npydir)
  - $(inputs.outdir)
  - $(inputs.outprefix)
  - $(inputs.npyconsolfile)
