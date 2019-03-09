$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Create SGLF (library) files
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 16
    ramMin: 100000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  bashscript:
    type: File
    label: Master script to create tile library (SGLF)
    default:
      class: File
      location: src/tilelib_chunk_v2CWL.sh
  tilepathmin:
    type: string
    label: Starting path in the tile library
  tilepathmax:
    type: string
    label: Last/Maximum path in the tile library
  fjcsv2sglf:
    type: string
    label: Tool to create tile library
    default: "/usr/local/bin/fjcsv2sglf"
  fjdir:
    type: Directory
    label: Directory of FastJ files
  fjt:
    type: string
    label: Tool to manipulate FastJ files
    default: "/usr/local/bin/fjt"
  tagset:
    type: File
    label: Compressed tagset in FASTA format
outputs:
  chunksglfs:
    type: File[]
    label: Output SGLF files
    outputBinding:
      glob: "lib/*sglf.gz*"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.tilepathmin)
  - $(inputs.tilepathmax)
  - $(inputs.fjcsv2sglf)
  - $(inputs.fjdir)
  - $(inputs.fjt)
  - $(inputs.tagset)
