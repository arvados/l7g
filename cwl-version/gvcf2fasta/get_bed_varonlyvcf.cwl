cwlVersion: v1.1
class: CommandLineTool
label: Get no call BED and variant only VCF from gVCF
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
    outdirMin: 40000
inputs:
  sampleid:
    type: string
    label: Sample ID
  vcf:
    type: File
    label: Input gVCF
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering  
  genomebed:
    type: File
    label: Whole genome BED
outputs:
  nocallbed:
    type: File
    label: No call BED of gVCF
    outputBinding:
      glob: "*_nocall.bed"
  varonlyvcf:
    type: File
    label: Variant only VCF
    outputBinding:
      glob: "*_varonly.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: /gvcf_regions/gvcf_regions.py
arguments:
  - prefix: "--min_GQ"
    valueFrom: $(inputs.gqcutoff)
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid).bed
  - shellQuote: False
    valueFrom: "&&"
  - "bedtools"
  - "subtract"
  - prefix: "-a"
    valueFrom: $(inputs.genomebed)
  - prefix: "-b"
    valueFrom: $(inputs.sampleid).bed
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid)_nocall.bed
  - shellQuote: False
    valueFrom: "&&"
  - "bgzip"
  - "-dc"
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "grep"
  - "-v"
  - "END="
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid)_varonly.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sampleid)_varonly.vcf.gz
