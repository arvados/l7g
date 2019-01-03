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
      location: ../../src/create-npyCWL.sh
    inputBinding:
      position: 1
  cgft:
    type: [File,string]
    default: "/usr/local/bin/cgft"
    inputBinding:
      position: 2
  cgfdirectory:
    type: Directory
    inputBinding:
      position: 3
  band2matrix:
    type: [File,string]
    default: "/usr/local/bin/band-to-matrix-npy"
    inputBinding:
      position: 4
  cnvrt2hiq:
    type: [File,string]
    default: "/usr/local/bin/npy-vec-to-hiq-1hot"
    inputBinding:
      position: 5
  makelist:
    type: File
    default: 
      class: File
      location: ../../src/create-list
    inputBinding:
      position: 6
  nthreads:
    type: string
    default: "16"
    inputBinding:
      position: 7
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "npy"
  out2:
    type: Directory
    outputBinding:
      glob: "npy-hiq"
