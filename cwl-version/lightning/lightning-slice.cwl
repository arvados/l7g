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
    arv:dockerCollectionPDH: 1f430e6dd9b6be0ae78d4cffde9b1fef+892
  ResourceRequirement:
    coresMin: 96
    ramMin: 660000
  arv:RuntimeConstraints:
    keep_cache: 6200
    outputDirType: keep_output_dir
inputs:
  datalibs:
    type:
      type: array
      items: File
  reflib:
    type: File
outputs:
  libdir:
    type: Directory
    outputBinding:
      glob: "."
baseCommand: [lightning, slice]
arguments:
  - "-local=true"
  - prefix: "-output-dir"
    valueFrom: $(runtime.outdir)
  - $(inputs.datalibs)
  - $(inputs.reflib)
