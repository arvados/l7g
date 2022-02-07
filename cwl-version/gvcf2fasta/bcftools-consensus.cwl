cwlVersion: v1.1
class: CommandLineTool
label: Convert VCF to FASTA with bcftools consensus
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  sampleid:
    type: string
    label: sample ID
  vcf:
    type: File
    label: Input VCF
    secondaryFiles: [.tbi]
  ref:
    type: File
    label: Reference FASTA
  haplotype:
    type: int
    label: Haplotype of sample (1 or 2)
  mask:
    type: File
    label: Mask BED region where FASTA sequence is filled with 'N'
outputs:
  fa:
    type: File
    label: Output FASTA 
    outputBinding:
      glob: "*fa.gz"
baseCommand: [bcftools, consensus]
arguments:
  - prefix: "--fasta-ref"
    valueFrom: $(inputs.ref)
  - prefix: "--haplotype"
    valueFrom: $(inputs.haplotype)
#  - prefix: "--mask"
#    valueFrom: $(inputs.mask)
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sampleid).$(inputs.haplotype).fa.gz
