cwlVersion: v1.1
class: Workflow
requirements:
  ScatterFeatureRequirement: {}

inputs:
  matchgenome_array:
    type: string[]
  libdir:
    type: Directory
  regions:
    type: File?
  threads:
    type: int
  mergeoutput:
    type: string
  expandregions:
    type: int

outputs:
  npydirs:
    type:
      type: array
      items: Directory
    outputSource: lightning-slice-numpy/npydir

steps:
  lightning-slice-numpy:
    run: lightning-slice-numpy.cwl
    scatter: matchgenome
    in:
      matchgenome: matchgenome_array
      libdir: libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
    out: [npydir]
