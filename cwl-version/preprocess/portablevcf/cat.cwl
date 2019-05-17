cwlVersion: v1.0
class: CommandLineTool
label: Concatenate files
hints:
  DockerRequirement:
    dockerPull: vcfutil
inputs:
  txts:
    type: File[]
    label: Text files
outputs:
  cattxt:
    type: stdout
    label: Concatenated text
baseCommand: cat
arguments:
  - $(inputs.txts)
stdout: catsummary.txt
