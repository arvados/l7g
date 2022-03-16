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
  phenotypesdir:
    type: Directory
  chrs: string[]
  snpeffdatadir: Directory
  genomeversion_array: string[]
  dbsnp:
    type: File
    secondaryFiles: [.csi]
  gnomaddir: Directory

outputs:
  stagednpydir:
    type: Directory[]
    outputSource: lightning-slice-numpy-anno2vcf-wf/stagednpydir
  stagedonehotnpydir:
    type: Directory[]
    outputSource: lightning-slice-numpy-anno2vcf-wf/stagedonehotnpydir
  stagedannotationdir:
    type:
      type: array
      items: [Directory, "null"]
    outputSource: lightning-slice-numpy-anno2vcf-wf/stagedannotationdir

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

  make-arrays:
    run: make-arrays.cwl
    in:
      matchgenome_array: matchgenome_array
      libdir_array: lightning-slice/libdir
      genomeversion_array: genomeversion_array
      regions_nestedarray: regions_nestedarray
      threads_array: threads_array
      mergeoutput_array: mergeoutput_array
      expandregions_array: expandregions_array
    out: [full_matchgenome_array, full_libdir_array, full_genomeversion_array, full_regions_array, full_threads_array, full_mergeoutput_array, full_expandregions_array, full_libname_array]

  lightning-slice-numpy-anno2vcf-wf:
    run: lightning-slice-numpy-anno2vcf-wf.cwl
    scatter: [matchgenome, libdir, genomeversion, regions, threads, mergeoutput, expandregions, libname]
    scatterMethod: dotproduct
    in:
      matchgenome: make-arrays/full_matchgenome_array
      libdir: make-arrays/full_libdir_array
      regions: make-arrays/full_regions_array
      threads: make-arrays/full_threads_array
      mergeoutput: make-arrays/full_mergeoutput_array
      expandregions: make-arrays/full_expandregions_array
      phenotypesdir: phenotypesdir
      libname: make-arrays/full_libname_array
      chrs: chrs
      snpeffdatadir: snpeffdatadir
      genomeversion: make-arrays/full_genomeversion_array
      dbsnp: dbsnp
      gnomaddir: gnomaddir
    out: [stagednpydir, stagedonehotnpydir, stagedannotationdir]
