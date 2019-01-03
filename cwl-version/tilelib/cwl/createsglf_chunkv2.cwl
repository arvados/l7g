$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Create SGLF (library) files
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    ramMin: 100000
    coresMin: 16
    coresMax: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Master script to create tile library (SGLF)
    inputBinding:
      position: 1
  tilepathmin:
    type: string
    label: Path to start at in the tile library
    inputBinding:
      position: 2
  tilepathmax:
    type: string
    label: Last/Maximum tile in the library path
    inputBinding:
      position: 3
  fjcsv2sglf:
    type: File
    label: Tool to create tile library
    inputBinding:
      position: 4
  datadir:
    type: Directory
    label: Directory of FastJ files
    inputBinding:
      position: 5
  fjt:
    type: File
    label: Tool to manipulate FastJ (text) files
    inputBinding:
      position: 6
  tagset:
    type: File
    label: Compressed tagset in FASTA format
    inputBinding:
      position: 7

outputs:
  out1:
    type: File[]
    label: Output for inclusion in tile library (SGLF)
    outputBinding:
      glob: "lib/*sglf.gz*"
