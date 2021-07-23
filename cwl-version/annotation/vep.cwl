cwlVersion: v1.1
class: CommandLineTool
label: Run Variant Effect Predictor given VCF
requirements:
  ShellCommandRequirement: {}
  DockerRequirement:
    dockerPull: ensemblorg/ensembl-vep:release_103
hints:
  ResourceRequirement:
    coresMin: 4
    ramMin: 20000
    tmpdirMin: 16000
inputs:
  vcf:
    type: File
    label: Input VCF
    secondaryFiles: [.tbi]
  sample:
    type: string
    label: Sample name
  vepcache:
    type: Directory
    label: Cache directory for Variant Effect Predictor
  assembly:
    type: string
    label: Assembly version
outputs:
  consequencevcf:
    type: File
    label: Consequence vcf
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: vep
arguments:
  - "--cache"
  - "--offline"
  - prefix: "--dir_cache"
    valueFrom: $(inputs.vepcache)
  - prefix: "--input_file"
    valueFrom: $(inputs.vcf)
  - "--check_existing"
  - prefix: "--assembly"
    valueFrom: $(inputs.assembly)
  - prefix: "--fork"
    valueFrom: "4"
  - "af_gnomad"
  - prefix: "--compress_output"
    valueFrom: "bgzip"
  - "--vcf"
  - prefix: "-o"
    valueFrom: $(inputs.sample)_csq.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_csq.vcf.gz
