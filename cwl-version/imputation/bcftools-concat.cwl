cwlVersion: v1.1
class: CommandLineTool
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: beagle5.4
  ResourceRequirement:
    coresMin: 2
    ramMin: 5000
    tmpdirMin: 10000
inputs:
  sample: string
  vcfs:
    type: File[]
    secondaryFiles: [.tbi]
outputs:
  vcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bcftools, concat]
arguments:
  - $(inputs.vcfs)
  - "-Oz"
  - prefix: "-o"
    valueFrom: $(inputs.sample)_rawimputed.vcf.gz
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_rawimputed.vcf.gz
