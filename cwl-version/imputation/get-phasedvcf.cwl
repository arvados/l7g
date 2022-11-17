cwlVersion: v1.1
class: CommandLineTool
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  sample: string
  vcf: File
outputs:
  phasedvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: zcat
arguments:
  - $(inputs.vcf)
  - shellQuote: false
    valueFrom: "|"
  - "egrep"
  - prefix: "-v"
    valueFrom: '0\|0|IMP'
  - shellQuote: false
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: false
    valueFrom: ">"
  - $(inputs.sample).vcf.gz
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample).vcf.gz
