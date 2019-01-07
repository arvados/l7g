$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 2
    ramMin: 13000
hints:
  arv:RuntimeConstraints:
    keep_cache: 1046
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
baseCommand: bash
inputs:
  bashscript:
    type: File
    default:
      class: File
      location: ../src/convertcgfCWL-empty-problem-tilepaths3.sh
    inputBinding:
      position: 1
  fjdir:
    type: Directory
    inputBinding:
      position: 2
  cgft:
    type: string
    default: "/usr/local/bin/cgft"
    inputBinding:
      position: 3
  fjt:
    type: string
    default: "/usr/local/bin/fjt"
    inputBinding:
      position: 4
  cglf:
    type: Directory
    inputBinding:
      position: 5
outputs:
  out1:
    type: File
    outputBinding:
      glob: "data/*.cgf"
