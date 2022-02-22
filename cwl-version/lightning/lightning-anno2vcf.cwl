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
    coresMin: 64
    ramMin: 500000
  arv:RuntimeConstraints:
    keep_cache: 83000
    outputDirType: keep_output_dir
inputs:
  annodir:
    type: Directory
  regions:
    type: File?
outputs:
  vcfdir:
    type: Directory
    outputBinding:
      glob: "."
baseCommand: [lightning, anno2vcf]
arguments:
  - "-local=true"
  - prefix: "-input-dir="
    valueFrom: $(inputs.annodir)
    separate: false
  - prefix: "-output-dir="
    valueFrom: $(runtime.outdir)
    separate: false
