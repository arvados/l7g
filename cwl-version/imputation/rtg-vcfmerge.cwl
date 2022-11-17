cwlVersion: v1.1
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  sample: string
  vcf:
    type: File
    secondaryFiles: [.tbi]
  phasedvcf:
    type: File
    secondaryFiles: [.tbi]
  imputedvcf:
    type: File
    secondaryFiles: [.tbi]
outputs:
  phasedimputedvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [rtg, vcfmerge]
arguments:
  - "--force-merge-all"
  - $(inputs.phasedvcf)
  - $(inputs.vcf)
  - $(inputs.imputedvcf)
  - prefix: "-o"
    valueFrom: $(inputs.sample)_phased_imputed.vcf.gz
