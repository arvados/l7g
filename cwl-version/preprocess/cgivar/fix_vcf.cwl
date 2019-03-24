cwlVersion: v1.0
class: CommandLineTool
requirements:
  InlineJavascriptRequirement: {}
inputs:
  fixscript: File
  vcf: File
outputs:
  fixedvcf: stdout
arguments:
  - $(inputs.fixscript)
  - $(inputs.vcf)
stdout: |
  ${
    return inputs.vcf.nameroot.split('.')[0]+".vcf"
  }
