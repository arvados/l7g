$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  refdirectory: Directory
  bashscript: File
  ref: string 
  reffa: File
  afn: File
  aidx: File
  tagdir: File
  l7g: File
  pasta: File
  refstream: File

outputs:
  out1:
    type: 
      type: array
      items:
        type: array
        items: Directory 
    outputSource: step2/out1
  out2:
    type:
      type: array 
      items:
        type: array
        items: File 
    outputSource: step2/out2

steps:
  step1:
    run: getfilesgff.cwl
    in: 
      refdirectory: refdirectory
    out: [out1]

  step2:
    scatter: gffInitial 
    scatterMethod: dotproduct
    in: 
      gffInitial: step1/out1
      bashscript: bashscript
      ref: ref
      reffa: reffa
      afn: afn
      aidx: aidx
      tagdir: tagdir
      l7g: l7g
      pasta: pasta
      refstream: refstream 
    run: convertgff.cwl
    out: [out1,out2]
