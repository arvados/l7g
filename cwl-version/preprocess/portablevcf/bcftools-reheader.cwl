cwlVersion: v1.0
class: CommandLineTool
label: Change the header of VCF
hints:
  DockerRequirement:
    dockerPull: vcfutil
inputs:
  header:
    type: File
    label: Header file
  vcfgz:
    type: File
    label: Input VCF
outputs:
  reheaderedvcfgz:
    type: File
    label: Reheadered VCF
    outputBinding:
      glob: "*vcf.gz"
baseCommand: [bcftools, reheader]
arguments:
  - prefix: "-h"
    valueFrom: $(inputs.header)
  - $(inputs.vcfgz)
  - prefix: "-o"
    valueFrom: $(inputs.vcfgz.basename)
