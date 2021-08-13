$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: lightning
  ResourceRequirement:
    coresMin: 96
    ramMin: 680000
  arv:RuntimeConstraints:
    keep_cache: 6500
inputs:
  lib1:
    type: File
  lib2:
    type: File
outputs:
  mergedlib:
    type: File
    outputBinding:
      glob: "*gob.gz"
baseCommand: [lightning, merge]
arguments:
  - "-local=true"
  - prefix: "-o"
    valueFrom: "library.gob.gz"
  - $(inputs.lib1)
  - $(inputs.lib2)
