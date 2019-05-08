cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
  - class: ResourceRequirement
    ramMin: 13000
inputs:
  vcf: File
  bed: File
outputs:
  sortedbed:
    type: File
    outputBinding:
      glob: "*.bed"