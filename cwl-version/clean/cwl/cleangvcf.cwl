$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: InlineJavascriptRequirement
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
  gvcfDir:
    type: Directory
    inputBinding:
      position: 2
  gvcfPrefix:
    type: string
    inputBinding:
      position: 3
  cleanvcf:
    type: [File,string]
    default: "/usr/local/bin/cleanvcf"
    inputBinding:
      position: 4
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "cleaned/*"
