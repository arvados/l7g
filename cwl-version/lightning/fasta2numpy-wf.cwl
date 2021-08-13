cwlVersion: v1.1
class: Workflow
requirements:
  StepInputExpressionRequirement: {}

inputs:
  tagset:
    type: File
  fastadir:
    type: Directory
  refdir:
    type: Directory
  chunks:
    type: int

outputs:
  mergedlib:
    type: File
  outdir:
    type: Directory

steps:
  lightning-import_data:
    run: lightning-import.cwl
    in:
      saveincomplete:
        valueFrom: "false"
      tagset: tagset
      fastadir: fastadir
    out: [stats, lib]
  lightning-import_ref:
    run: lightning-import.cwl
    in:
      saveincomplete:
        valueFrom: "true"
      tagset: tagset
      fastadir: refdir
    out: [stats, lib]
  lightning-merge:
    run: lightning-merge.cwl
    in:
      lib1: lightning-import_data/lib
      lib2: lightning-import_ref/lib
    out: [mergedlib]
  lightning-export-numpy:
    run: lightning-export-numpy.cwl
    in:
      lib: lightning-merge/mergedlib
      chunks: chunks
    out: [outdir]
