cwlVersion: v1.0
class: CommandLineTool
requirements:
  class: ShellCommandRequirement: {}
  class: ResourceRequirement
    #coresMin: 2
    ramMin: 13000
inputs:
  vcf: File
  bed: File
outputs:
  vcfgz:
    type: File
    outputBinding:
      glob: "*.vcf.gz" 
    secondaryFiles: [.tbi]
baseCommand: [gunzip]
arguments:
    valueFrom: $(inputs.vcf.basename).gz
  - shellQuote: false
    valueFrom: "&&"
  - "bedtools"
  - "intersect"
  - "-header" 
  - prefix: "-a" 
    valueFrom: $(inputs.vcf)
  - prefix: "-b" 
    valueFrom: $(inputs.bed)
  - prefix: "-f" 
    valueFrom: "1" 
  - shellQuote: false
    valueFrom: "|" 
  - "bgzip" 
  - "-c" 
  - shellQuote: false
    valueFrom: ">" 
  - $(inputs.vcf.basename).gz
  - shellQuote: false
    valueFrom: "&&" 
  - "tabix" 
  - prefix: "-p" 
    valueFrom: "vcf" 
  - $(inputs.vcf.basename).gz

  #TODO - add a preliminary step of `gunzip` on the .vcf.gz