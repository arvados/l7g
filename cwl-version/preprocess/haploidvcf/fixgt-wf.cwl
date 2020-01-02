cwlVersion: v1.0
class: Workflow
label: Scatter to fix VCF by changing haploid calls and processing chrM
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
    outputSource: fixgt/fixedvcf

steps:
  getfiles:
    run: getfiles.cwl
    in:
      dir: vcfdir
    out: [vcfs]
  fixgt:
    run: fixgt.cwl
    scatter: [vcf]
    in:
      vcf: getfiles/vcfs
    out: [fixedvcf]
