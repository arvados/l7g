$namespaces:
  arv: "http://arvados.org/cwl#"
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
  refdir:
    type: Directory
  batchsize:
    type: int
  regions:
    type: File?
  matchgenome:
    type: string
  threads:
    type: int
  mergeoutput:
    type: string
  expandregions:
    type: int

outputs: []

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
    in:
      saveincomplete:
        valueFrom: "true"
      tagset: tagset
      fastadirs: refdir
    out: [lib]

  lightning-slice:
    run: lightning-slice.cwl
    in:
      datalibs: lightning-import_data/lib
      reflib: lightning-import_refs/lib
    out: [libdir]

  lightning-tiling-stats:
    run: lightning-tiling-stats.cwl
    in:
      libdir: lightning-slice/libdir
    out: [bed]

  lightning-slice-numpy:
    run: lightning-slice-numpy.cwl
    in:
      matchgenome: matchgenome
      libdir: lightning-slice/libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
    out: [outdir, npys, samplescsv, chunktagoffsetcsv]
