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
  pcanpy: File
  pcasamplescsv: File
  phenotypesdir: Directory
  xcomponent: string
  ycomponent: string
outputs:
  png:
    type: File
    outputBinding:
      glob: "*.png"
baseCommand: [lightning, plot]
arguments:
  - "-local=true"
  - prefix: "-i="
    valueFrom: $(inputs.pcanpy)
    separate: false
  - prefix: "-o="
    valueFrom: "plot_$(inputs.xcomponent)-$(inputs.ycomponent).png"
    separate: false
  - prefix: "-samples="
    valueFrom: $(inputs.pcasamplescsv)
    separate: false
  - prefix: "-phenotype="
    valueFrom: $(inputs.phenotypesdir)
    separate: false
  - "-phenotype-cat1-column=7"
  - prefix: "-x="
    valueFrom: $(inputs.xcomponent)
    separate: false
  - prefix: "-y="
    valueFrom: $(inputs.ycomponent)
    separate: false
