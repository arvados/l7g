$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Sort BED by natural ordering (1,2,10,M,X)
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
inputs:
  bed:
   type: File
   label: BED to be sorted by natural ordering
outputs:
  sortedbed:
    type: File
    label: BED sorted by natural ordering
    outputBinding:
      glob: "*.bed"
baseCommand: sort
arguments:
  - prefix: "-k1,1V"
    valueFrom: "-k2,2n"
  - $(inputs.bed)
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.bed.basename)