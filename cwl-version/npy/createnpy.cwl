$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Create NumPy vectors from cgfs by tile path
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 16
    ramMin: 100000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  bashscript:
    type: File
    label: Master script for creating the NumPy arrays
    default:
      class: File
      location: src/create-npyCWL.sh
  cgft:
    type: string
    label: Compact genome format tool
    default: "/usr/local/bin/cgft"
  cgfdir:
    type: Directory
    label: Directory of compact genome format files
  band2matrix:
    type: string
    label: Tool to convert band (path) information into NumPy array
    default: "/usr/local/bin/band-to-matrix-npy"
  cnvrt2hiq:
    type: string
    label: Tool to create NumPy files for high quality arrays
    default: "/usr/local/bin/npy-vec-to-hiq-1hot"
  makelist:
    type: File
    label: Tool for saving dataset names
    default:
      class: File
      location: src/create-list
outputs:
  npydir:
    type: Directory
    label: Directory of NumPy arrays
    outputBinding:
      glob: "npy"
  npyhiqdir:
    type: Directory
    label: Directory of high quality NumPy arrays
    outputBinding:
      glob: "npy-hiq"
  names:
    type: File
    label: File listing sample names
    outputBinding:
      glob: "npy/names"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.cgft)
  - $(inputs.cgfdir)
  - $(inputs.band2matrix)
  - $(inputs.cnvrt2hiq)
  - $(inputs.makelist)
  - $(runtime.cores)
