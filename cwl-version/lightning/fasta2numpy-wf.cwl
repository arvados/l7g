cwlVersion: v1.2
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  tagset:
    type: File
  fastadirs:
    type:
      type: array
      items: Directory
  refdirs:
    type:
      type: array
      items: Directory
  batchsize:
    type: int
  matchgenome_array:
    type: string[]
  regions_nestedarray:
    type:
      type: array
      items:
        type: array
        items: [File, "null"]
  threads_array:
    type: int[]
  mergeoutput_array:
    type: string[]
  expandregions_array:
    type: int[]

outputs:
  libdirs:
    type:
      type: array
      items: Directory
    outputSource: lightning-slice/libdir
  npydirs:
    type:
      type: array
      items: Directory
    outputSource: scatter3-lightning-slice-numpy-anno2vcf-wf/npydirs
  vcfdirs:
    type:
      type: array
      items: Directory
    outputSource: scatter3-lightning-slice-numpy-anno2vcf-wf/vcfdirs

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

  lightning-import_refs:
    run: lightning-import.cwl
    scatter: fastadirs
    in:
      saveincomplete:
        valueFrom: "true"
      tagset: tagset
      fastadirs: refdirs
    out: [lib]

  lightning-slice:
    run: lightning-slice.cwl
    scatter: reflib
    in:
      datalibs: lightning-import_data/lib
      reflib: lightning-import_refs/lib
    out: [libdir]

  scatter3-lightning-slice-numpy-anno2vcf-wf:
    run: scatter3-lightning-slice-numpy-anno2vcf-wf.cwl
    in:
      matchgenome_array: matchgenome_array
      libdir_array: lightning-slice/libdir
      regions_nestedarray: regions_nestedarray
      threads_array: threads_array
      mergeoutput_array: mergeoutput_array
      expandregions_array: expandregions_array
    out: [npydirs, vcfdirs]
