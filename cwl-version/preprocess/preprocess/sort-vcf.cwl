cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: fbh/vcfpreprocess
inputs: 
  vcf:
    type: File
    inputBinding:
      position: 1
      prefix: -c
outputs: 
  sortedvcf:
    type: stdout
baseCommand: vcf-sort
stdout: "sorted-output.vcf"