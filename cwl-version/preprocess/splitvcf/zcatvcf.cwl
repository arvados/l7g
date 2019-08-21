$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Concatenate VCFs split by chromosomes
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  bashscript:
    type: File
    label: Master script to concatenate VCFs
    default:
      class: File
      location: src/zcatvcf.sh
  vcfdir:
    type: Directory
    label: Input VCFs directory
outputs:
  vcf:
    type: File
    label: Concatenated VCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.vcfdir)
