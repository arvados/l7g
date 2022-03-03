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
  trimmedvcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
baseCommand: awk
arguments:
  - '{if ($1 ~ /^#/ || $4 != $5) print $0}'
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sample).vcf.gz
