cwlVersion: v1.0
class: CommandLineTool
label: Get unzipped sglf size and find paths above a given threshold
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
hints:
  ResourceRequirement:
    coresMin: 8
inputs:
  bashscript:
    type: File
    label: Master script to get unzipped sglf size
    default:
      class: File
      location: src/getsglfsize.sh
  lib:
    type: Directory
    label: Tile library
  threshold:
    type: int
    label: Threshold for unzipped sglf size in MiB
outputs:
  sglfsize:
    type: File
    label: Unzipped sglf size
    outputBinding:
      glob: "*tsv"
  skippaths:
    type: File
    label: Paths above the given threshold
    outputBinding:
      glob: "*txt"
arguments:
  - $(inputs.bashscript)
  - $(inputs.lib)
  - $(inputs.threshold)
