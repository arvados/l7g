$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Convert gVCF to FastJ
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 13000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  bashscript:
    type: File
    label: Master script to create a FastJ for each gVCF
    default:
      class: File
      location: src/convert2fastjCWL
  gvcf:
    type: File
    label: Input gVCF
    secondaryFiles: [.tbi]
  ref:
    type: string
    label: Reference genome
  reffa:
    type: File
    label: Reference genome in FASTA format
  afn:
    type: File
    label: Compressed assembly fixed width file
  aidx:
    type: File
    label: Assembly index file
  refM:
    type: string
    label: Mitochondrial reference genome
  reffaM:
    type: File
    label: Reference mitochondrial genome in FASTA format
  afnM:
    type: File
    label: Compressed mitochondrial assembly fixed width file
  aidxM:
    type: File
    label: Mitochondrial assembly index file
  seqidM:
    type: string
    label: Mitochondrial naming scheme for gVCF
  tagset:
    type: File
    label: Compressed tagset in FASTA format
  l7g:
    type: string
    label: Lightning application for parsing and searching assembly files
    default: "/usr/local/bin/l7g"
  pasta:
    type: string
    label: Tool for handling verbose stream oriented format
    default: "/usr/local/bin/pasta"
  refstream:
    type: string
    label: Tool for streaming and converting variant call formats
    default: "/usr/local/bin/refstream"
  tile_assembly:
    type: string
    label: Tool to extract information from the tile assembly files
    default: "/usr/local/bin/tile-assembly"
outputs:
  fjdir:
    type: Directory
    label: Directories of FastJs
    outputBinding:
      glob: "stage/*"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.gvcf)
  - $(inputs.ref)
  - $(inputs.reffa)
  - $(inputs.afn)
  - $(inputs.aidx)
  - $(inputs.refM)
  - $(inputs.reffaM)
  - $(inputs.afnM)
  - $(inputs.aidxM)
  - $(inputs.seqidM)
  - $(inputs.tagset)
  - $(inputs.l7g)
  - $(inputs.pasta)
  - $(inputs.refstream)
  - $(inputs.tile_assembly)
