$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Create SGLF (library) files
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 16
    coresMax: 16
    ramMin: 100000
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
    default:
      class: File
      location: ../src/tilelib_chunk_v2CWL.sh
    inputBinding:
      position: 1
  tilepathmin:
    type: string
    label: Starting path in the tile library
    inputBinding:
      position: 2
  tilepathmax:
    type: string
    label: Last/Maximum path in the tile library
    inputBinding:
      position: 3
  fjcsv2sglf:
    type: string
    label: Tool to create tile library
    default: "/usr/local/bin/fjcsv2sglf"
    inputBinding:
      position: 4
  datadir:
    type: Directory
    label: Directory of FastJ files
    inputBinding:
      position: 5
  fjt:
    type: string
    label: Tool to manipulate FastJ files
    default: "/usr/local/bin/fjt"
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
    label: Output SGLF files
    outputBinding:
      glob: "lib/*sglf.gz*"
