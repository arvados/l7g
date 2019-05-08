cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
inputs:
  vcf: File
  bed: File
outputs:
  sortedbed:
    type: File
    outputBinding:
      glob: "*.bed"
baseCommand: sorted
  - valueFrom: "-k1,1V"
  - "-k2,2n"
  - $(inputs.bed)