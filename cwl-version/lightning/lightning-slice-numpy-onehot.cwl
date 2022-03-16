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
    ramMin: 660000
  arv:RuntimeConstraints:
    keep_cache: 83000
    outputDirType: keep_output_dir
inputs:
  matchgenome:
    type: string
  libdir:
    type: Directory
  regions:
    type: File?
  threads:
    type: int
  mergeoutput:
    type: string
  expandregions:
    type: int
  phenotypesdir:
    type: Directory
outputs:
  outdir:
    type: Directory
    outputBinding:
      glob: "."
  onehotcolumnsnpy:
    type: File
    outputBinding:
      glob: "onehot-columns.npy"
  onehotnpy:
    type: File
    outputBinding:
      glob: "onehot.npy"
  csv:
    type: File
    outputBinding:
      glob: "samples.csv"
baseCommand: [lightning, slice-numpy]
arguments:
  - "-local=true"
  - prefix: "-input-dir="
    valueFrom: $(inputs.libdir)
    separate: false
  - prefix: "-output-dir="
    valueFrom: $(runtime.outdir)
    separate: false
  - prefix: "-match-genome="
    valueFrom: $(inputs.matchgenome)
    separate: false
  - prefix: "-regions="
    valueFrom: $(inputs.regions)
    separate: false
  - prefix: "-threads="
    valueFrom: $(inputs.threads)
    separate: false
  - prefix: "-merge-output="
    valueFrom: $(inputs.mergeoutput)
    separate: false
  - prefix: "-expand-regions="
    valueFrom: $(inputs.expandregions)
    separate: false
  - "-chunked-onehot=true"
  - "-single-onehot=true"
  - prefix: "-chi2-case-control-file="
    valueFrom: $(inputs.phenotypesdir)
    separate: false
  - "-chi2-case-control-column=AD"
  - "-chi2-p-value=0.01"
  - "-min-coverage=0.9"
