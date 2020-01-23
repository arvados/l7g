cwlVersion: v1.0
class: Workflow
label: Scatter to fix VCF by processing chrM
requirements:
  ScatterFeatureRequirement: {}
inputs:
  vcfdir:
    type: Directory
    label: Input VCF directory
  filterjs:
    type: File
    label: Javascript code for filtering

outputs:
  fixedvcfs:
    type: File[]
    label: Fixed VCFs
    outputSource: fixchrm/fixedvcf

steps:
  getfiles:
    run: getfiles.cwl
    in:
      dir: vcfdir
    out: [vcfs]
  fixchrm:
    run: fixchrm.cwl
    scatter: [vcf]
    in:
      vcf: getfiles/vcfs
      filterjs: filterjs
    out: [fixedvcf]
