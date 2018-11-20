$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
    ramMin: 13000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    inputBinding:
      position: 1
  gffDir:
    type: Directory
    inputBinding:
      position: 2
  gffPrefix:
    type: string
    inputBinding:
      position: 3
  filter_gvcf: 
    type: File
    inputBinding:
      position: 4
  cutoff:
    type: string
    inputBinding:
      position: 5
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "filtered/*"
