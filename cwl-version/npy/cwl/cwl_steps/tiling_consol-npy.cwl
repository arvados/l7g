$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Consolidates multiple NumPy arrays into one large array
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
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
    label: Master script to consolidate tile path NumPy arrays
    inputBinding:
      position: 1
    default:
      class: File
      location: ../../src/allconsolCWL.sh
  indir:
    type: Directory
    label: Name of input directory
    inputBinding:
      position: 2
  outdir:
    type: string?
    label: Name of output directory
    inputBinding:
      position: 3
    default: "outdir"
  outprefix:
    type: string?
    label: Prefix for consolidated arrays
    inputBinding:
      position: 4
    default: "all"
  npyconsolfile:
    type: File?
    label: Program to consolidated NumPy arrays
    inputBinding:
      position: 5
    default:
      class: File
      location: ../../src/buildArvados/dest/npy-consolidate
outputs:
  out1:
    type: Directory
    label: Consolidated NumPy arrays
    outputBinding:
      glob: $(inputs.outdir)
