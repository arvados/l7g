$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Process and create cgf files from FastJ files.
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    ramMin: 10000
    coresMin: 2
hints:
  arv:RuntimeConstraints:
    keep_cache: 1046
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Main bash script for converting FastJ to cgf
    inputBinding:
      position: 1
  fjdir:
    type: Directory
    label: FastJ directory
    inputBinding:
      position: 2
  cgft:
    type: File
    label: Compact Genome Format Tool, a tool to manipulate and inspect cgf files
    inputBinding:
      position: 3
  fjt:
    type: File
    label: a tool to manipulate FastJ (text) files.
    inputBinding:
      position: 4
  cglf:
    type: Directory
    label: Tile library location
    inputBinding:
      position: 5
outputs:
  out1:
    type: File
    outputBinding:
      glob: "data/*.cgf"
