cwlVersion: v1.0
class: CommandLineTool
label: Intersect VCF with BED
requirements:
  ShellCommandRequirement: {}
inputs:
  vcf:
    type: File
    label: Input VCF
  bed:
    type: File
    label: Input BED
outputs:
  vcfgz:
    type: File
    label: Output VCF with records inside the BED region
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bedtools, intersect]
arguments:
  - "-header"
  - prefix: "-a"
    valueFrom: $(inputs.vcf)
  - prefix: "-b"
    valueFrom: $(inputs.bed)
  - prefix: "-f"
    valueFrom: "1"
  - shellQuote: false
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: false
    valueFrom: ">"
  - $(inputs.vcf.basename).gz
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - prefix: "-p"
    valueFrom: "vcf"
  - $(inputs.vcf.basename).gz
