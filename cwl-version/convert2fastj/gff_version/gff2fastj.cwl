$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
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
    default:
      class: File
      location: src/convert2fastjCWL
  gff:
    type: File
  ref:
    type: string
  reffa:
    type: File
  afn:
    type: File
  aidx:
    type: File
  refM:
    type: string
  reffaM:
    type: File
  afnM:
    type: File
  aidxM:
    type: File
  seqidM:
    type: string 
  tagset:
    type: File
  l7g:
    type: string
    default: "/usr/local/bin/l7g"
  pasta:
    type: string
    default: "/usr/local/bin/pasta"
  refstream:
    type: string
    default: "/usr/local/bin/refstream"
outputs:
  fjdir:
    type: Directory
    outputBinding:
      glob: "stage/*"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.gff)
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
