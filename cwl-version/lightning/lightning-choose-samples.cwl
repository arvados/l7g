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
    ramMin: 10000
  arv:RuntimeConstraints:
    keep_cache: 20000
    outputDirType: keep_output_dir
inputs:
  matchgenome: string
  libdir: Directory
  phenotypesdir: Directory
  trainingsetsize: float
  randomseed: int
outputs:
  samplescsv:
    type: File
    outputBinding:
      glob: "samples.csv"
baseCommand: [lightning, choose-samples]
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
  - prefix: "-case-control-file="
    valueFrom: $(inputs.phenotypesdir)
    separate: false
  - "-case-control-column=AD"
  - prefix: "-training-set-size="
    valueFrom: $(inputs.trainingsetsize)
    separate: false
  - prefix: "-random-seed="
    valueFrom: $(inputs.randomseed)
    separate: false
