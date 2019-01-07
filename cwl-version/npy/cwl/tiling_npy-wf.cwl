$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g 

inputs:
  cgfdirectory: Directory
  nthreads:
    type: string
    default: "16"
  outdir:
    type: string
    default: "outdir"
  outprefix:
    type: string
    default: "all"

outputs:
  out1:
    type: Directory
    outputSource: step2/out1

steps:
  step1:
    run: tiling_create-npy.cwl
    in:
      cgfdirectory: cgfdirectory
      nthreads: nthreads
    out: [out1,out2]

  step2:
    run: tiling_consol-npy.cwl
    in: 
      indir: step1/out1
      outdir: outdir
      outprefix: outprefix
    out: [out1]
