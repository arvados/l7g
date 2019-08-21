cwlVersion: v1.0
class: CommandLineTool
label: Check tile library for correct formatting and spurious characters
$namespaces:
  arv: "http://arvados.org/cwl#"
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4000
inputs:
  script:
    type: File
    label: Master script to run tile library check
    default:
      class: File
      location: src/sglf-sanity-check
  sglfdir:
    type: Directory
    label: Tile library directory
outputs:
  log:
    type: File
    label: Validation logs
    outputBinding:
      glob: "*log"
baseCommand: bash
arguments:
  - $(inputs.script)
  - $(inputs.sglfdir)
  - "sglf.log"
