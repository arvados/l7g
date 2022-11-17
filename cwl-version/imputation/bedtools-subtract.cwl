cwlVersion: v1.1
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  sample: string
  a: File
  b: File
outputs:
  subtractbed: stdout
baseCommand: [bedtools, subtract]
arguments:
  - prefix: "-a"
    valueFrom: $(inputs.a)
  - prefix: "-b"
    valueFrom: $(inputs.b)
stdout: $(inputs.sample)_subtract.bed
