cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement 
    dockerPull: fbh/vcfpreprocess
  - class: ShellCommandRequirement  
inputs: 
  sortedvcf:
    type: File
    inputBinding:
      prefix: -a
      position: 1
  sortedbed:
    type: File
    inputBinding:
      prefix: -b
      position: 2
  intersectoverlap:
    type: string
    default: "-F 1"
    inputBinding:
      shellQuote: False 
      position: 3
# Print the header from the A file prior to results      
  headeroption:
    type: string
    default: "-header"
    inputBinding:
      shellQuote: False
      position: 4
outputs: 
  intersectedvcf:
    type: stdout
baseCommand: intersectBed
stdout: "intersected-output.vcf"