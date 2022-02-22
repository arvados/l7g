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
  sampleid:
    type: string
  vcfdir:
    type: Directory
  bashscript:
    type: File
    default:
      class: File
      location: src/concatenate.sh
outputs:
  vcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
arguments:
  - $(inputs.bashscript)
  - $(inputs.vcfdir)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid).vcf.gz
