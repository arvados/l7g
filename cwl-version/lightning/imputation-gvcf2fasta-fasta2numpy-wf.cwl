$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.2
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}

inputs:
  sampleids:
    type: string[]
  splitvcfdirs:
    type: Directory[]
  gqcutoff:
    type: int
  genomebed:
    type: File
  ref:
    type: File
  chrs: string[]
  refsdir: Directory
  mapsdir: Directory
  panelnocallbed: File
  panelcallbed: File
  tagset:
    type: File
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
  scatter-gvcf2fasta_splitvcf-imputation-wf:
    run: ../gvcf2fasta/scatter-gvcf2fasta_splitvcf-imputation-wf.cwl
    in:
      sampleids: sampleids
      splitvcfdirs: splitvcfdirs
      gqcutoff: gqcutoff
      genomebed: genomebed
      ref: ref
      chrs: chrs
      refsdir: refsdir
      mapsdir: mapsdir
      panelnocallbed: panelnocallbed
      panelcallbed: panelcallbed
    out: [fas]

  make-fastadirs:
    run: make-fastadirs.cwl
    in:
      fas: scatter-gvcf2fasta_splitvcf-imputation-wf/fas
    out: [fastadirs]

  batch-dirs:
    run: batch-dirs.cwl
    in:
      dirs: make-fastadirs/fastadirs
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
    out: [outdir, npys, chunktagoffsetcsv]
