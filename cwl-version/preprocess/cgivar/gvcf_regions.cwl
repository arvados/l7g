cwlVersion: v1.0
class: CommandLineTool
inputs:
  vcf: File
outputs:
  bed: stdout
baseCommand: /gvcf_regions/gvcf_regions.py
arguments:
  - $(inputs.vcf)
  - "--unreported_is_called"
stdout: $(inputs.vcf.nameroot).bed
