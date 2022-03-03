cwlVersion: v1.1
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  counts: File[]
  bashscript:
    type: File
    default:
      class: File
      location: src/totalcounts.sh
outputs:
  summary:
    type: stdout
arguments:
  - $(inputs.bashscript)
  - $(inputs.counts)
stdout: summary.txt
