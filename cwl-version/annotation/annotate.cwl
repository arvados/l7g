$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.1
class: CommandLineTool
label: Get HGVS annotation for a given path
requirements:
  DockerRequirement:
    dockerPull: tileannotation
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4000
inputs:
  bashscript:
    type: File
    label: Master script to get HGVS
    default:
      class: File
      location: src/annotate_hgvs.sh
  get_hgvs:
    type: File
    label: Python script to get HGVS
    default:
      class: File
      location: ../../tools/get_hgvs/get_hgvs.py
  pathstr:
    type: string
    label: Input path string
  ref:
    type: File
    label: Reference genome FASTA
  tilelib:
    type: Directory
    label: Input tile library
  varnum:
    type: int
    label: The number of tile variants to be annotated per step
  assembly:
    type: File
    label: Compressed assembly fixed width file
    secondaryFiles: [^.fwi, .gzi]
outputs:
  annotation:
    type: File
    label: HGVS annotation in csv format
    outputBinding:
      glob: "*csv"
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.get_hgvs)
  - $(inputs.pathstr)
  - $(inputs.ref)
  - $(inputs.tilelib)
  - $(inputs.varnum)
  - $(inputs.assembly)
