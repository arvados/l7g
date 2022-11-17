cwlVersion: v1.1
class: CommandLineTool
label: Convert VCF to FASTA with bcftools consensus
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
  mask:
    type: File
    label: Mask BED region where FASTA sequence is filled with 'N'
  bashscript:
    type: File
    label: Script to run bcftools consensus
    default:
      class: File
      location: src/bcftools-consensus.sh
outputs:
  fas:
    type: File[]
    label: Output FASTAs 
    outputBinding:
      glob: "*fa.gz"
arguments:
  - $(inputs.bashscript)
  - $(inputs.sampleid)
  - $(inputs.vcf)
  - $(inputs.ref)
  - $(inputs.mask)
