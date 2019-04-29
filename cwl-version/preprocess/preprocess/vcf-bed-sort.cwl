$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"

cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: fbh/vcfpreprocess
  - class: ShellCommandRequirement
  - class: ResourceRequirement
    coresMin: 1
    ramMin: 13000
inputs:
  vcf: File
  bed: File
outputs:
  sortedvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz" 
    secondaryFiles: [.tbi]
  sortedbed:
    type: File
    outputBinding:
      glob: "*.bed"
baseCommand: vcf-sort
arguments:
  - prefix: "-c"
    valueFrom: $(inputs.vcf)
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.vcf.basename)
  - shellQuote: False
    valueFrom: "&&"
  - "sort"
  - "-k1,1V"
  - "-k2,2n" 
  - $(inputs.bed)
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.bed.basename)
  - shellQuote: False
    valueFrom: "&&"
  - "bedtools"
  - "intersect"
  - "-header" 
  - prefix: "-a" 
    valueFrom: $(inputs.vcf.basename)
  - prefix: "-b" 
    valueFrom: $(inputs.bed.basename)
  - prefix: "-f" 
    valueFrom: "1" 
  - shellQuote: false
    valueFrom: "|" 
  - "bgzip" 
  - shellQuote: false
    valueFrom: ">" 
  - $(inputs.vcf.basename).gz
  - shellQuote: false
    valueFrom: "&&" 
  - "tabix" 
  # - prefix: "-p" 
  #   valueFrom: "vcf" 
  - $(inputs.vcf.basename).gz