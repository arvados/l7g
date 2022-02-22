cwlVersion: v1.2
class: Workflow
requirements:
  InlineJavascriptRequirement: {}

inputs:
  matchgenome:
    type: string
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
  npydir:
    type: Directory
    outputSource: lightning-slice-numpy/npydir
  vcfdir:
    type: Directory?
    outputSource: lightning-anno2vcf/vcfdir

steps:
  lightning-slice-numpy:
    run: lightning-slice-numpy.cwl
    in:
      matchgenome: matchgenome
      libdir: libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
    out: [npydir]

  lightning-anno2vcf:
    run: lightning-anno2vcf.cwl
    when: $(inputs.regions == null)
    in:
      annodir: lightning-slice-numpy/npydir
      regions: regions
    out: [vcfdir]
    
