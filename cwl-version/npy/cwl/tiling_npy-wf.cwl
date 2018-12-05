$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create numpy arrays from cgf, merge all numpy arrays into one array
doc: |
    Merge individual Lightning numpy arrays broken out by tilepath into a single numpy matrix
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  bashscriptmain_create:
    type: File?
    label: Master bash script for creating the numpy arrays
  bashscriptmain_consol:
    type: File?
    label: Bash script to Consolidate individual Lightning numpy arrays broken out by tilepath into a single numpy matrix
  cgft:
    type: ["null", "File", "string"]
    label: compact genome format tool
  cgfdirectory:
    type: Directory
    label: Directory for compact genome format files
  band2matrix:
    type: File?
    label: Compiled C++ band-to-matrix-npy convert band information into a Lightning tile numpy array
  cnvrt2hiq:
    type: File?
    label: Compiled C++ npy-vec-to-hiq-1hot create 'flat' numpy hiq tile vector arrays and its info file
  makelist:
    type: File?
    label: used for saving the names of the datasets as a numpy array
  nthreads:
    type: string?
    label: Number of threads to use
  outdir:
    type: string?
    label: Name of output directory
  outprefix:
    type: string?
    label: Prefix or patch to prepend to output Directory
  npyconsolfile:
    type: File?
    label: Name of consolidated numpy array

outputs:
  out1:
    type: Directory
    outputSource: step2/out1

steps:
  step1:
    run: ../cwl/cwl_steps/tiling_create-npy.cwl
    in:
      bashscriptmain: bashscriptmain_create
      cgft: cgft
      cgfdirectory: cgfdirectory
      band2matrix: band2matrix
      cnvrt2hiq: cnvrt2hiq
      makelist: makelist
      nthreads: nthreads
    out: [out1,out2]
    label: Output numpy arrays

  step2:
    run: ../cwl/cwl_steps/tiling_consol-npy.cwl
    in:
      bashscriptmain: bashscriptmain_consol
      indir: step1/out1
      outdir: outdir
      outprefix: outprefix
      npyconsolfile: npyconsolfile
    out: [out1]
    label: Output consolidated numpy arrays
