$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create NumPy arrays by tile path from cgfs, merge all NumPy arrays into single array
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g 

inputs:
  cgfdirectory:
    type: Directory
    label: Directory of compact genome format files

outputs:
  out1:
    type: Directory
    label: Output consolidated NumPy arrays
    outputSource: step2/out1
  names:
    type: File
    label: File listing sample names
    outputSource: step1/names

steps:
  step1:
    run: tiling_create-npy.cwl
    in:
      cgfdirectory: cgfdirectory
    out: [out1,out2,names]

  step2:
    run: tiling_consol-npy.cwl
    in: 
      indir: step1/out1
    out: [out1]
