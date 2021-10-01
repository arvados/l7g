$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: lightning
  ResourceRequirement:
    coresMin: 96
    ramMin: 640000
  arv:RuntimeConstraints:
    keep_cache: 6500
inputs:
  saveincomplete:
    type: string
  tagset:
    type: File
  fastadir:
    type: Directory
outputs:
  stats:
    type: File
    outputBinding:
      glob: "*json"
  lib:
    type: File
    outputBinding:
      glob: "*gob.gz"
baseCommand: [lightning, import]
arguments:
  - "-local=true"
  - "-skip-ooo=true"
  - "-output-tiles=true"
  - "-batches=1"
  - "-batch=0"
  - prefix: "-save-incomplete-tiles="
    valueFrom: $(inputs.saveincomplete)
    separate: false
  - prefix: "-match-chromosome"
    valueFrom: "^(chr)?([0-9]+|X|Y|MT?)$"
  - prefix: "-output-stats"
    valueFrom: "stats.json"
  - prefix: "-tag-library"
    valueFrom: $(inputs.tagset)
  - prefix: "-o"
    valueFrom: "library.gob.gz"
  - $(inputs.fastadir)
