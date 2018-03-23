cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
  - class: SubworkflowFeatureRequirement
  - class: ScatterFeatureRequirement

inputs:
  script: File
  tilepaths: string[]
  fastj_base_dir: Directory
  tagset: File
  outdir: string

outputs:
  result:
    type: Directory
    outputSource: gather/result

steps:
  process:
    run: sglf-create-tilepath.cwl
    scatter: tilepath
    in:
      script: script
      tilepath: tilepaths
      fastj_base_dir: fastj_base_dir
      tagset: tagset
      outdir: outdir
    out: [result]
  gather:
    run: sglf-create-gather.cwl
    in:
      idirs: process/result
    out: [result]
