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
      location: ../../src/create-npyCWL.sh
  cgft:
    type: File? 
    inputBinding:
      position: 2
    default:
      class: File
      location: ../../src/buildArvados/dest/cgft
  cgfdirectory:
    type: Directory
    inputBinding:
      position: 3
  band2matrix:
    type: File?
    inputBinding:
      position: 4 
    default:
      class: File
      location: ../../src/buildArvados/dest/band-to-matrix-npy
  cnvrt2hiq:
    type: File?
    inputBinding:
      position: 5
    default:
      class: File
      location: ../../src/buildArvados/dest/npy-vec-to-hiq-1hot
  makelist:
    type: File?
    inputBinding:
      position: 6
    default: 
      class: File
      location: ../../src/create-list
  nthreads:
    type: string?
    inputBinding:
      position: 7
    default: "16"
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "npy"
  out2:
    type: Directory
    outputBinding:
      glob: "npy-hiq"
