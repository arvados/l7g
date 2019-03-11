$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
inputs:
  bashscript:
    type: File
    default:
      class: File
      location: src/clean-gff-header.sh
  gff:
    type: File
outputs:
  cleangff:
    type: File
    outputBinding:
      glob: "*gz"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.gff)
