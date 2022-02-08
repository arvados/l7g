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
  vcf:
    type: File
    secondaryFiles: [.tbi]
  genomebed:
    type: File
  nocallbed:
    type: File
outputs:
  intersectvcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
baseCommand: [bedtools, subtract]
arguments:
  - prefix: "-a"
    valueFrom: $(inputs.genomebed)
  - prefix: "-b"
    valueFrom: $(inputs.nocallbed)
  - shellQuote: False
    valueFrom: "|"
  - "bedtools"
  - "intersect"
  - prefix: "-a"
    valueFrom: $(inputs.vcf)
  - prefix: "-b"
    valueFrom: "stdin"
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid).vcf.gz
