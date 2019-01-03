$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Filters gVCFs by a specified quality cutoff
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Master script to control filtering
    inputBinding:
      position: 1
  gffDir:
    type: Directory
    label: Directory of gVCF files
    inputBinding:
      position: 2
  gffPrefix:
    type: string
    label: Prefix for gVCF files
    inputBinding:
      position: 3
  filter_gvcf:
    type: File
    label: Code that filters gVCFs
    inputBinding:
      position: 4
  cutoff:
    type: string
    label: Filtering cutoff threshold
    inputBinding:
      position: 5
outputs:
  out1:
    type: Directory
    label: Directory of Filtered gVCFs
    outputBinding:
      glob: "filtered/*"
