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
  excludebed: File
outputs:
  filteredvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [rtg, vcffilter]
arguments:
  - prefix: "-i"
    valueFrom: $(inputs.vcf)
  - prefix: "-o"
    valueFrom: $(inputs.sample).vcf.gz
  - prefix: "--exclude-bed"
    valueFrom: $(inputs.excludebed)
