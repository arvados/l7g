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
    coresMin: 2
    ramMin: 8000
  arv:RuntimeConstraints:
    keep_cache: 10000
    outputDirType: keep_output_dir
inputs:
  libdir: Directory
  regions: File?
outputs:
  bed:
    type: File
    outputBinding:
      glob: "*bed"
baseCommand: [lightning, tiling-stats]
arguments:
  - "-local=true"
  - prefix: "-input-dir"
    valueFrom: $(inputs.libdir)
  - prefix: "-output-dir"
    valueFrom: $(runtime.outdir)
