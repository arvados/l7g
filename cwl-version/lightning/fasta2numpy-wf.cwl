cwlVersion: v1.1
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  MultipleInputFeatureRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  tagset:
    type: File
  fastadirs:
    type:
      type: array
      items: Directory
  refdir:
    type: Directory
  batchsize:
    type: int

outputs:
  libdir:
    type: Directory
    outputSource: lightning-slice/libdir
  npydir:
    type: Directory
    outputSource: lightning-slice-numpy/npydir

steps:
  batch-dirs:
    run: batch-dirs.cwl
    in:
      dirs: fastadirs
      batchsize: batchsize
    out: [batches]

  lightning-import_data:
    run: lightning-import.cwl
    scatter: fastadirs
    in:
      saveincomplete:
        valueFrom: "false"
      tagset: tagset
      fastadirs: batch-dirs/batches
    out: [lib]

  lightning-import_ref:
    run: lightning-import.cwl
    in:
      saveincomplete:
        valueFrom: "true"
      tagset: tagset
      fastadirs: refdir
    out: [lib]

  lightning-slice:
    run: lightning-slice.cwl
    in:
      libs:
        source: [lightning-import_data/lib, lightning-import_ref/lib]
        linkMerge: merge_flattened
    out: [libdir]

  lightning-slice-numpy:
    run: lightning-slice-numpy.cwl
    in:
      libdir: lightning-slice/libdir
    out: [npydir]
