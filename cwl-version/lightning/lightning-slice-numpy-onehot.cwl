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
  matchgenome: string
  libdir: Directory
  regions: File?
  threads: int
  mergeoutput: string
  expandregions: int
  phenotypesdir: Directory
outputs:
  outdir:
    type: Directory
    outputBinding:
      glob: "."
  npys:
    type: File[]
    outputBinding:
      glob: "*npy"
  samplescsv:
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
  - "-single-onehot=true"
  - prefix: "-chi2-case-control-file="
    valueFrom: $(inputs.phenotypesdir)
    separate: false
  - "-chi2-case-control-column=AD"
  - "-chi2-p-value=0.01"
  - "-min-coverage=0.9"
