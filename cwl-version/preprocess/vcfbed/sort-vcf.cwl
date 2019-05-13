$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Sort VCF by natural ordering (1,2,10,M,X)
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
inputs:
  vcf:
    type: File
    label: Compressed VCF to be sorted by natural ordering
outputs:
  sortedvcf:
    type: File
    label: Compressed VCF sorted by natural ordering
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: vcf-sort
arguments:
  - prefix: "-c"
    valueFrom: $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.vcf.basename)
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - prefix: "-p"
    valueFrom: "vcf"
  - $(inputs.vcf.basename)