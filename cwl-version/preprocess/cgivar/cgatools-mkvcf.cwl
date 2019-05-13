cwlVersion: v1.0
class: CommandLineTool
label: Convert CGIVAR to VCF
inputs:
  cgascript:
    type: File
    label: Script invoking cgatools
  reference:
    type: File
    label: CRR reference used for cgatools
  cgivar:
    type: File
    label: Input CGIVAR
  sample:
    type: string
    label: Sample name
outputs:
  vcf:
    type: stdout
    label: Output VCF
arguments:
  - $(inputs.cgascript)
  - $(inputs.reference)
  - $(inputs.cgivar)
stdout: $(inputs.sample).vcf
