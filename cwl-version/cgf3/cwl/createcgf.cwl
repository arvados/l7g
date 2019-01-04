$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Process and create cgf files from FastJ files
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
    label: Master script to convert FastJs to cgfs
    inputBinding:
      position: 1
  fjdir:
    type: Directory
    label: Input directory of FastJs
    inputBinding:
      position: 2
  cgft:
    type: File
    label: Tool to manipulate and inspect cgf files
    inputBinding:
      position: 3
  fjt:
    type: File
    label: Tool to manipulate FastJ files
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
    label: Output cgfs
    outputBinding:
      glob: "data/*.cgf"
