$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.2
class: CommandLineTool
requirements:
  NetworkAccess:
    networkAccess: true
hints:
  DockerRequirement:
    dockerPull: lightning
  ResourceRequirement:
    coresMin: 96
    ramMin: 670000
  arv:RuntimeConstraints:
    keep_cache: 6200
    outputDirType: keep_output_dir
inputs:
  saveincomplete:
    type: string
  tagset:
    type: File
  fastadirs:
    type:
      - Directory
      - type: array
        items: Directory
outputs:
  lib:
    type: File
    outputBinding:
      glob: "*gob.gz"
baseCommand: [lightning, import]
arguments:
  - "-local=true"
  - "-loglevel=info"
  - "-skip-ooo=true"
  - "-output-tiles=true"
  - "-batches=1"
  - "-batch=0"
  - prefix: "-save-incomplete-tiles="
    valueFrom: $(inputs.saveincomplete)
    separate: false
  - prefix: "-match-chromosome"
    valueFrom: "^(chr)?([0-9]+|X|Y|M)$"
  - prefix: "-output-stats"
    valueFrom: "stats.json"
  - prefix: "-tag-library"
    valueFrom: $(inputs.tagset)
  - prefix: "-o"
    valueFrom: "library.gob.gz"
  - $(inputs.fastadirs)
