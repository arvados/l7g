cwlVersion: v1.0
class: CommandLineTool
inputs:
  cgascript: File
  reference: File
  cgivar: File
  sample: string
outputs:
  vcf: stdout
arguments:
  - $(inputs.cgascript)
  - $(inputs.reference)
  - $(inputs.cgivar)
stdout: $(inputs.sample).vcf
