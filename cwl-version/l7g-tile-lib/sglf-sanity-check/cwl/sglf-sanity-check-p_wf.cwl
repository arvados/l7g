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
  idir: Directory
  nBatch: string

outputs:
  result:
    type: Directory
    outputSource: gather/result

steps:
  batch_input:
    run: sglf-sanity-check-p_batch-input.cwl
    in:
      idir: idir
      nBatch: nBatch
    out: [batchDir, batchName]
  process:
    run: sglf-sanity-check-p.cwl
    scatter: [ idir, outFileName ]
    scatterMethod: dotproduct
    in:
      script: script
      idir: batch_input/batchDir
      outFileName: batch_input/batchName
      ncore: nBatch
    out: [result]
  gather:
    run: sglf-sanity-check-p_gather.cwl
    in:
      idirs: process/result
    out: [result]

