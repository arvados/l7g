cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: fbh/preprocess-vcfbed
  - class: ResourceRequirement
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
  - shellQuote: False
  - prefix: "-a" 
    valueFrom: $(inputs.vcf)
  - prefix: "-b" 
    valueFrom: $(inputs.bed)  
  - prefix: "-f" 
    valueFrom: "1" 
  - "|" 
  - "bgzip"
  - shellQuote: False
    prefix: "-c"
    valueFrom: ">"
  - $(inputs.vcf.basename)
  - shellQuote: False
    valueFrom: "&&"
  - "tabix" 
  - prefix: "-p" 
    valueFrom: "vcf" 
  - $(inputs.vcf.basename)
  - shellQuote: False
