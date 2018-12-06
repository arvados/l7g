$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
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
    label: Bash script that iterates over the FastJ to create paths
    inputBinding:
      position: 1
  tilepathmin:
    type: string
    label: Beginning tile library path [0]
    inputBinding:
      position: 2
  tilepathmax:
    type: string
    label: Last/Maximum tile library path
    inputBinding:
      position: 3
  fjcsv2sglf:
    type: File
    label: Compiled C++ that creates 2bit sequence and tile ID and size
    inputBinding:
      position: 4
  datadir:
    type: Directory
    label: Directory in Keep for Data
    inputBinding:
      position: 5
  fjt:
    type: File
    label: fjt is a tool to manipulate FastJ (text) files
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
    outputBinding:
      glob: "lib/*sglf.gz*"
