$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  bashscriptmain: File
  cgft: File
  cgfdirectory: Directory
  band2matrix: File
  cnvrt2hiq: File
  makelist: File
  nthreads: string
  outdir: string
  outprefix: string
  npyconsolfile: File

outputs:
  out1:
    type: Directory[]
    outputSource: step2/out1

steps:
  step1:
    run: /cwl_steps/tiling_create-npy.cwl
    in:
      datafilenames: datafilenames
      refdirectory: refdirectory
    out: [out1,out2]

  step2:
    run: /cwl_steps/
    in:
      bashscript: bashscript
      gvcfDir: step1/out1
      gvcfPrefix: step1/out2
      cleanvcf: cleanvcf
    out: [out1]
                                                              46,1          Bot

