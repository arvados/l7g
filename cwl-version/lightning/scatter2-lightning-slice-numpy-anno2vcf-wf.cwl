cwlVersion: v1.2
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  SubworkflowFeatureRequirement: {}

inputs:
  matchgenome_array:
    type: string[]
  libdir:
    type: Directory
  regions_array:
    type:
      type: array
      items: [File, "null"]
  threads_array:
    type: int[]
  mergeoutput_array:
    type: string[]
  expandregions_array:
    type: int[]

outputs:
  npydirs:
    type:
      type: array
      items: Directory
    outputSource: flatten-array_npydirs/flattenedarray
  vcfdirs:
    type:
      type: array
      items: Directory
    outputSource: flatten-array_vcfdirs/flattenedarray

steps:
  scatter1-lightning-slice-numpy-anno2vcf-wf:
    run: scatter1-lightning-slice-numpy-anno2vcf-wf.cwl
    scatter: [regions, threads, mergeoutput, expandregions]
    scatterMethod: dotproduct
    in:
      matchgenome_array: matchgenome_array
      libdir: libdir
      regions: regions_array
      threads: threads_array
      mergeoutput: mergeoutput_array
      expandregions: expandregions_array
    out: [npydirs, vcfdirs]

  flatten-array_npydirs:
    run: flatten-array.cwl
    in:
      nestedarray: scatter1-lightning-slice-numpy-anno2vcf-wf/npydirs
    out: [flattenedarray]

  flatten-array_vcfdirs:
    run: flatten-array.cwl
    in:
      nestedarray: scatter1-lightning-slice-numpy-anno2vcf-wf/vcfdirs
    out: [flattenedarray]
