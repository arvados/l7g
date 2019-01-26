$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 1
    coresMax: 1
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Master script to control cleaning
    default:
      class: File
      location: ../src/cleanCWL.sh
    inputBinding:
      position: 1
  gvcfDir:
    type: Directory
    label: Directory with gVCF files
    inputBinding:
      position: 2
  gvcfPrefix:
    type: string
    label: Prefix of gVCF files
    inputBinding:
      position: 3
  cleanvcf:
    type: string
    label: Tool to clean gVCFs
    default: "/usr/local/bin/cleanvcf"
    inputBinding:
      position: 4
outputs:
  out1:
    type: Directory
    label: Directory of clean gVCFs
    outputBinding:
      glob: "cleaned/*"
