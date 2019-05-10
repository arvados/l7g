cwlVersion: v1.0
class: CommandLineTool
label: Sort BED by natural ordering (1,2,10,MT,X)
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
inputs:
  bed: File
outputs:
  sortedbed:
    type: File
    outputBinding:
      glob: "*.bed"
baseCommand: sort
arguments:
  - prefix: "-k1,1V"
    valueFrom: "-k2,2n"
  - $(inputs.bed)
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.bed.basename)