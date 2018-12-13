$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Create NumPy vectors based on tile library paths
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
    label: Main workflow script
    inputBinding:
      position: 1
    default:
      class: File
      location: ../../src/create-npyCWL.sh
  cgft:
    type: ["null",File,string]
    label: Compact genome format tool
    inputBinding:
      position: 2
    default: "usr/bin/cgft"
  cgfdirectory:
    type: Directory
    label: Directory for compact genome format files
    inputBinding:
      position: 3
  band2matrix:
    type: File?
    label: Convert band information into a Lightning tile NumPy array.
    inputBinding:
      position: 4
    default:
      class: File
      location: ../../src/buildArvados/dest/band-to-matrix-npy
  cnvrt2hiq:
    type: File?
    label: Create flat NumPy hiq tile vector arrays and its info file
    inputBinding:
      position: 5
    default:
      class: File
      location: ../../src/buildArvados/dest/npy-vec-to-hiq-1hot
  makelist:
    type: File?
    label: Save the names of the datasets as a NumPy array
    inputBinding:
      position: 6
    default:
      class: File
      location: ../../src/create-list
  nthreads:
    type: string?
    label: Number of threads to use
    inputBinding:
      position: 7
    default: "16"
outputs:
  out1:
    type: Directory
    label: Directory for NumPy arrays
    outputBinding:
      glob: "npy"
  out2:
    type: Directory
    label: Directory for hiq NumPy arrays
    outputBinding:
      glob: "npy-hiq"
