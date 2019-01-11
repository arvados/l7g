cwlVersion: v1.0
class: Workflow
label: Workflow to create a FASTA file for the tagset
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

inputs:
  cytobandFn:
    type: File
  bigwigFn:
    type: File
  refFaFn:
    type: File
  script:
    type: File

outputs:
  result:
    type: File[]
    outputSource: tagset/result

steps:
  tagset:
    run: tagset.cwl
    in:
      cytobandFn: cytobandFn
      bigwigFn: bigwigFn
      refFaFn: refFaFn
      script: script
    out: [result]
