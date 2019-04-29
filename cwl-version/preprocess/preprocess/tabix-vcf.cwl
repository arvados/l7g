cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: fbh/vcfpreprocess
inputs:
  bgzipvcf:
    type: File
    inputBinding:
      position: 1
outputs:
  tabixvcf:
    type: File
baseCommand: tabix