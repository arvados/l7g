w $namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create NumPy arrays by tile path from cgfs, merge all NumPy arrays into single array
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
    label: Script to consolidate tile path arrays into a single NumPy matrix
  cgft:
    type: ["null", "File", "string"]
    label: Compact genome format tool
  cgfdirectory:
    type: Directory
    label: Directory of compact genome format files
  band2matrix:
    type: File?
    label: Tool to convert band (path) information into NumPy array
  cnvrt2hiq:
    type: File?
    label: Tool to create NumPy files for high quality arrays
  makelist:
    type: File?
    label: Tool for saving dataset names
  nthreads:
    type: string?
    label: Number of threads to use
  outdir:
    type: string?
    label: Name of output directory
  outprefix:
    type: string?
    label: Prefix for consolidated arrays
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
