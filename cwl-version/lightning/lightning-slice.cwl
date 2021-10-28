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
    coresMin: 96
    ramMin: 660000
  arv:RuntimeConstraints:
    keep_cache: 6200
    outputDirType: keep_output_dir
inputs:
  libs:
    type:
      type: array
      items: File
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
  - $(inputs.libs)
