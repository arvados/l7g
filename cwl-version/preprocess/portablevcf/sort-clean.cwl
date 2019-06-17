cwlVersion: v1.0
class: CommandLineTool
label: Sort VCF and clean duplicate calls
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
inputs:
  cleanvcf:
    type: File
    label: Code that cleans VCFs
  vcfgz:
    type: File
    label: Input VCF
outputs:
  cleanvcfgz:
    type: File
    label: Clean VCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: vcf-sort
arguments:
  - "-c"
  - $(inputs.vcfgz)
  - shellQuote: False
    valueFrom: "|"
  - $(inputs.cleanvcf)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.vcfgz.basename)
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - prefix: "-p"
    valueFrom: "vcf"
  - $(inputs.vcfgz.basename)
