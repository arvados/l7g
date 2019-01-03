$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g 

inputs:
  bashscriptmain_create:
    type: File
    default:
      class: File
      location: ../src/create-npyCWL.sh
  bashscriptmain_consol:
    type: File
    default:
      class: File
      location: ../src/allconsolCWL.sh
  cgft:
    type: [File,string]
    default: "/usr/local/bin/cgft"
  cgfdirectory: Directory
  band2matrix:
    type: [File,string]
    default: "/usr/local/bin/band-to-matrix-npy"
  cnvrt2hiq:
    type: [File,string]
    default: "/usr/local/bin/npy-vec-to-hiq-1hot"
  makelist:
    type: File
    default: 
      class: File
      location: ../src/create-list
  nthreads:
    type: string
    default: "16"
  outdir:
    type: string
    default: "outdir"
  outprefix:
    type: string
    default: "all"
  npyconsolfile:
    type: [File,string]
    default: "/usr/local/bin/npy-consolidate"

outputs:
  out1:
    type: Directory
    outputSource: step2/out1

steps:
  step1:
    run: cwl_steps/tiling_create-npy.cwl
    in:
      bashscriptmain: bashscriptmain_create
      cgft: cgft
      cgfdirectory: cgfdirectory
      band2matrix: band2matrix
      cnvrt2hiq: cnvrt2hiq
      makelist: makelist
      nthreads: nthreads
    out: [out1,out2]

  step2:
    run: cwl_steps/tiling_consol-npy.cwl
    in: 
      bashscriptmain: bashscriptmain_consol
      indir: step1/out1
      outdir: outdir
      outprefix: outprefix
      npyconsolfile: npyconsolfile
    out: [out1]
