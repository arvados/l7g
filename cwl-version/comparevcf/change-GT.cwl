cwlVersion: v1.2
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  sampleid: string
  suffix: string
  vcf: File
  header: File
  bashscript: File
outputs:
  modifiedvcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
arguments:
  - $(inputs.bashscript)
  - $(inputs.sampleid)
  - $(inputs.suffix)
  - $(inputs.vcf)
  - $(inputs.header)
