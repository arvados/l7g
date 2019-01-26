$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Create NumPy vectors from cgfs by tile path
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
    label: Master script for creating the NumPy arrays
    default:
      class: File
      location: ../src/create-npyCWL.sh
    inputBinding:
      position: 1
  cgft:
    type: string
    label: Compact genome format tool
    default: "/usr/local/bin/cgft"
    inputBinding:
      position: 2
  cgfdirectory:
    type: Directory
    label: Directory of compact genome format files
    inputBinding:
      position: 3
  band2matrix:
    type: string
    label: Tool to convert band (path) information into NumPy array
    default: "/usr/local/bin/band-to-matrix-npy"
    inputBinding:
      position: 4
  cnvrt2hiq:
    type: string
    label: Tool to create NumPy files for high quality arrays
    default: "/usr/local/bin/npy-vec-to-hiq-1hot"
    inputBinding:
      position: 5
  makelist:
    type: File
    label: Tool for saving dataset names
    default:
      class: File
      location: ../src/create-list
    inputBinding:
      position: 6
  nthreads:
    type: string
    label: Number of threads to use
    default: "16"
    inputBinding:
      position: 7
outputs:
  out1:
    type: Directory
    label: Directory of NumPy arrays
    outputBinding:
      glob: "npy"
  out2:
    type: Directory
    label: Directory of high quality NumPy arrays
    outputBinding:
      glob: "npy-hiq"
  names:
    type: File
    label: File listing sample names
    outputBinding:
      glob: "npy/names"
