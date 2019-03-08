$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Process and create cgf files from FastJ files
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 13000
hints:
  arv:RuntimeConstraints:
    keep_cache: 1046
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  bashscript:
    type: File
    label: Master script to convert FastJs to cgfs
    default:
      class: File
      location: src/convertcgfCWL-empty-problem-tilepaths3.sh
  fjdir:
    type: Directory
    label: Input directory of FastJs
  cgft:
    type: string
    label: Tool to manipulate and inspect cgf files
    default: "/usr/local/bin/cgft"
  fjt:
    type: string
    label: Tool to manipulate FastJ files
    default: "/usr/local/bin/fjt"
  lib:
    type: Directory
    label: Tile library directory
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.fjdir)
  - $(inputs.cgft)
  - $(inputs.fjt)
  - $(inputs.lib)
outputs:
  cgf:
    type: File
    label: Output cgf
    outputBinding:
      glob: "data/*.cgf"
