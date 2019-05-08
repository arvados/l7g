cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
inputs:
  vcf: File
  bed: File
outputs:
  sortedvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: vcf-sort
arguments:
  - prefix: "-c"
    valueFrom: $(inputs.vcf)
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.vcf.basename)
  - shellQuote: False
    valueFrom: "&&"