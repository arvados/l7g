$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create numpy arrays from cgf and merge all numpy arrays into a single array
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g 
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  bashscriptmain_create: File?
  bashscriptmain_consol: File?
  cgft: ["null", "File", "string"]
  cgfdirectory: Directory
  band2matrix: File?
  cnvrt2hiq: File?
  makelist: File?
  nthreads: string?
  outdir: string?
  outprefix: string?
  npyconsolfile: File?

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

  step2:
    run: ../cwl/cwl_steps/tiling_consol-npy.cwl
    in: 
      bashscriptmain: bashscriptmain_consol
      indir: step1/out1
      outdir: outdir
      outprefix: outprefix
      npyconsolfile: npyconsolfile
    out: [out1]
