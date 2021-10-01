$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: lightning
  ResourceRequirement:
    coresMin: 32
    ramMin: 720000
  arv:RuntimeConstraints:
    keep_cache: 2000
inputs:
  lib:
    type: File
  chunks:
    type: int
outputs:
  outdir:
    type: Directory
    outputBinding:
      glob: "."
baseCommand: [lightning, export-numpy]
arguments:
  - "-local=true"
  - "-one-hot=false"
  - prefix: "-input-dir"
    valueFrom: $(inputs.lib)
  - prefix: "-output-dir"
    valueFrom: $(runtime.outdir)
  - prefix: "-output-annotations"
    valueFrom: "annotations.csv"
  - prefix: "-output-onehot2tilevar"
    valueFrom: "onehot2tilevar.csv"
  - prefix: "-output-labels"
    valueFrom: "labels.csv"
  - prefix: "-expand-regions"
    valueFrom: "0"
  - prefix: "-chunks"
    valueFrom: $(inputs.chunks)
  - "-max-variants=-1"
  - "-min-coverage=0.000000"
  - "-max-tag=-1"
