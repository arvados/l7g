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
    outputSource: lightning-slice-numpy-anno2vcf-wf/npydir
  vcfdirs:
    type:
      type: array
      items:
        - "null"
        - Directory
    outputSource: lightning-slice-numpy-anno2vcf-wf/vcfdir

steps:
  lightning-slice-numpy-anno2vcf-wf:
    run: lightning-slice-numpy-anno2vcf-wf.cwl
    scatter: matchgenome
    in:
      matchgenome: matchgenome_array
      libdir: libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
    out: [npydir, vcfdir]
