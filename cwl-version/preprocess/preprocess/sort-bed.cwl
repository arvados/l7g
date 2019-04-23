cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement 
    dockerPull: fbh/vcfpreprocess
  - class: ShellCommandRequirement  
inputs: 
  unsortedbed:
    type: File
    inputBinding:
      position: 1
  sort-function:
    type: string
    default: "-k1,1V -k2,2n"
    inputBinding:
      shellQuote: False 
      position: 2
outputs: 
  sortedbed:
    type: stdout
baseCommand: sort
stdout: "sorted-output.bed"