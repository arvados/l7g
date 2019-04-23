cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: fbh/vcfpreprocess
inputs:
  sortedvcf:
    type: File
    inputBinding:
      position: 1
outputs:
  bgzipvcf:
    type: File
baseCommand: bgzip