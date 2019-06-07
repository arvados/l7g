cwlVersion: v1.0
class: CommandLineTool
label: Fix VCF with an extra period in the INFO field
requirements:
  InlineJavascriptRequirement: {}
inputs:
  fixscript:
    type: File
    label: Script to fix VCF
  vcf:
    type: File
    label: Input VCF
outputs:
  fixedvcf:
    type: stdout
    label: Fixed VCF
arguments:
  - $(inputs.fixscript)
  - $(inputs.vcf)
stdout: $(inputs.vcf.nameroot).vcf
