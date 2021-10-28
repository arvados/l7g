$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: CommandLineTool
requirements:
  NetworkAccess:
    networkAccess: true
hints:
  DockerRequirement:
    dockerPull: lightning
  ResourceRequirement:
    coresMin: 64
    ramMin: 660000
  arv:RuntimeConstraints:
    keep_cache: 83000
    outputDirType: keep_output_dir
inputs:
  libdir:
    type: Directory
outputs:
  npydir:
    type: Directory
    outputBinding:
      glob: "."
baseCommand: [lightning, slice-numpy]
arguments:
  - "-local=true"
  - prefix: "-input-dir"
    valueFrom: $(inputs.libdir)
  - prefix: "-output-dir"
    valueFrom: $(runtime.outdir)
  - prefix: "-threads"
    valueFrom: "10"
  - prefix: "-expand-regions"
    valueFrom: "0"
  - "-max-variants=-1"
  - "-min-coverage=0.000000"
  - "-max-tag=-1"
