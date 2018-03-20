cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  script: File
  gvcfFn: File
  tagset: File
  tileassembly: File
  refFaFn: File
  refName: string
  outName: string

outputs:
  result:
    type: Directory
    outputSource: hupgp_gvcf_to_fastj/result

steps:
  hupgp_gvcf_to_fastj:
    run: gvcf-to-fastj.cwl
    in:
      script: script
      gvcfFn: gvcfFn
      tagset: tagset
      tileassembly: tileassembly
      refFaFn: refFaFn
      refName: refName
      outName: outName
    out: [result]
