$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.1
class: CommandLineTool
label: Get HGVS annotation for a given path
requirements:
  DockerRequirement:
    dockerPull: diff-fasta
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4000
inputs:
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
  assembly:
    type: File
    label: Compressed assembly fixed width file
    secondaryFiles: [^.fwi, .gzi]
  bashscript:
    type: File
    label: Bashscript for diff-fasta
    default:
      class: File
      location: src/diff-fasta.sh
  samplefastalimit:
    type: int
    label: Sample FASTA limit
    default: 5000
outputs:
  annotation:
    type: stdout
    label: HGVS annotation in csv format
arguments:
  - $(inputs.get_hgvs)
  - $(inputs.pathstr)
  - $(inputs.ref)
  - $(inputs.tilelib)
  - $(inputs.assembly)
  - $(inputs.bashscript)
  - prefix: "--samplefastalimit"
    valueFrom: $(inputs.samplefastalimit)
stdout: $(inputs.pathstr).csv
