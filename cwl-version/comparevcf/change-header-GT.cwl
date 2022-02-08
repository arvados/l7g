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
  suffix:
    type: string
  header:
    type: File
  vcf:
    type: File
  bashscript:
    type: File
    default:
      class: File
      location: src/change-header-GT.sh
outputs:
  modifiedvcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
arguments:
  - $(inputs.bashscript)
  - $(inputs.header)
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid)_$(inputs.suffix).vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sampleid)_$(inputs.suffix).vcf.gz
