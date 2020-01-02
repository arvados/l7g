cwlVersion: v1.0
class: Workflow
label: Scatter to fix VCF by removing GT fields that point to <NON_REF> and processing chrM
requirements:
  ScatterFeatureRequirement: {}
inputs:
  vcfdir:
    type: Directory
    label: Input VCF directory

outputs:
  fixedvcfs:
    type: File[]
    label: Fixed VCFs
    outputSource: fixnonref/fixedvcf

steps:
  getfiles:
    run: getfiles.cwl
    in:
      dir: vcfdir
    out: [vcfs]
  fixnonref:
    run: fixnonref.cwl
    scatter: [vcf]
    in:
      vcf: getfiles/vcfs
    out: [fixedvcf]
