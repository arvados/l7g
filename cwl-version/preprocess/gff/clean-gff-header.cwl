cwlVersion: v1.0
class: CommandLineTool
label: Clean GFF to remove header lines
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
inputs:
  bashscript:
    type: File
    label: Master script to clean GFF
    default:
      class: File
      location: src/clean-gff-header.sh
  gff:
    type: File
    label: Input GFF
outputs:
  cleangff:
    type: File
    label: Clean GFF
    outputBinding:
      glob: "*gz"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.gff)
