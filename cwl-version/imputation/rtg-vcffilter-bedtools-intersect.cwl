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
  vcf:
    type: File
    secondaryFiles: [.tbi]
  bed: File
outputs:
  filteredvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [rtg, vcffilter]
arguments:
  - "--remove-overlapping"
  - prefix: "-i"
    valueFrom: $(inputs.vcf)
  - prefix: "-o"
    valueFrom: "-"
  - shellQuote: false
    valueFrom: "|"
  - "bedtools"
  - "intersect"
  - "-header"
  - prefix: "-f"
    valueFrom: "1"
  - prefix: "-a"
    valueFrom: "stdin"
  - prefix: "-b"
    valueFrom: $(inputs.bed)
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
