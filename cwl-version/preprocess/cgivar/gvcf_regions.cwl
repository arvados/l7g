cwlVersion: v1.0
class: CommandLineTool
label: Extract called region from GVCF
inputs:
  vcf:
    type: File
    label: Input GVCF
outputs:
  bed:
    type: stdout
    label: BED region of GVCF
baseCommand: /gvcf_regions/gvcf_regions.py
arguments:
  - $(inputs.vcf)
  - "--unreported_is_called"
stdout: $(inputs.vcf.nameroot).bed
