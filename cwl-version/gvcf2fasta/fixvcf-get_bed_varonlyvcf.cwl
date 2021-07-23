cwlVersion: v1.1
class: CommandLineTool
label: Fix gVCF and get no call BED and variant only VCF from gVCF
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 8000
inputs:
  bashscript:
    type: File
    label: Bash script
    default:
      class: File
      location: src/fixvcf-get_bed_varsonlyvcf.sh
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
arguments:
  - $(inputs.bashscript)
  - $(inputs.sampleid)
  - $(inputs.vcf)
  - $(inputs.gqcutoff)
  - $(inputs.genomebed)
