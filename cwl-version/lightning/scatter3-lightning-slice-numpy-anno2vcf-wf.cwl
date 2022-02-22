cwlVersion: v1.2
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
  SubworkflowFeatureRequirement: {}

inputs:
  matchgenome_array:
    type: string[]
  libdir_array:
    type: Directory[]
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
  scatter2-lightning-slice-numpy-anno2vcf-wf:
    run: scatter2-lightning-slice-numpy-anno2vcf-wf.cwl
    scatter: [libdir, regions_array]
    scatterMethod: dotproduct
    in:
      matchgenome_array: matchgenome_array
      libdir: libdir_array
      regions_array: regions_nestedarray
      threads_array: threads_array
      mergeoutput_array: mergeoutput_array
      expandregions_array: expandregions_array
    out: [npydirs, vcfdirs]

  flatten-array_npydirs:
    run: flatten-array.cwl
    in:
      nestedarray: scatter2-lightning-slice-numpy-anno2vcf-wf/npydirs
    out: [flattenedarray]

  flatten-array_vcfdirs:
    run: flatten-array.cwl
    in:
      nestedarray: scatter2-lightning-slice-numpy-anno2vcf-wf/vcfdirs
    out: [flattenedarray]
