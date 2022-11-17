cwlVersion: v1.1
class: CommandLineTool
label: Concatenate and get no call BED and variant only VCF from gVCF
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
    outdirMin: 40000
inputs:
  sampleid:
    type: string
    label: Sample ID
  splitvcfdir:
    type: Directory
    label: Input directory of split gVCFs
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering  
  genomebed:
    type: File
    label: Whole genome BED
  bashscript:
    type: File
    label: Script to untar and concatenate vcf tar ball
    default:
      class: File
      location: src/concat-get_bed_varonlyvcf.sh
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
  - $(inputs.splitvcfdir)
  - $(inputs.gqcutoff)
  - $(inputs.genomebed)
