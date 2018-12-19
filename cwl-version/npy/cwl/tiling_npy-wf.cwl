w $namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create NumPy arrays from cgf, merge all NumPy arrays into one array by tile path
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  bashscriptmain_create:
    type: File?
    label: Master script for creating the NumPy arrays
  bashscriptmain_consol:
    type: File?
    label: Script to consolidate tile path NumPy arrays into a single NumPy matrix
  cgft:
    type: ["null", "File", "string"]
    label: Compact genome format tool
  cgfdirectory:
    type: Directory
    label: Directory for compact genome format (cgf) files
  band2matrix:
    type: File?
    label: Tool to convert band information into a Lightning tile NumPy array
  cnvrt2hiq:
    type: File?
    label: Tool to create numpy files for high quality tiles
  makelist:
    type: File?
    label: Used for saving the names of the datasets as a NumPy array
  nthreads:
    type: string?
    label: Number of threads to use
  outdir:
    type: string?
    label: Name of output directory
  outprefix:
    type: string?
    label: Prefix to prepend to consolidated NumPy arrays
  npyconsolfile:
    type: File?
    label: Program to consolidated NumPy arrays

outputs:
  out1:
    type: Directory
    outputSource: step2/out1
    label: Output consolidated NumPy arrays

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

  step2:
    run: ../cwl/cwl_steps/tiling_consol-npy.cwl
    in:
      bashscriptmain: bashscriptmain_consol
      indir: step1/out1
      outdir: outdir
      outprefix: outprefix
      npyconsolfile: npyconsolfile
    out: [out1]
