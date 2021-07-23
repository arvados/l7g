cwlVersion: v1.1
class: CommandLineTool
label: Run Variant Effect Predictor given VCF
requirements:
  DockerRequirement:
    dockerPull: ensemblorg/ensembl-vep:release_103
  ResourceRequirement:
    ramMin: 8000
    tmpdirMin: 16000
inputs:
  vcf:
    type: File
    label: Input VCF
  vepcache:
    type: Directory
    label: Cache directory for Variant Effect Predictor
  assembly:
    type: string
    label: Assembly version
outputs:
  consequence:
    type: stdout
    label: Consequence of the mutation
baseCommand: vep
arguments:
  - "--cache"
  - "--offline"
  - prefix: "--dir_cache"
    valueFrom: $(inputs.vepcache)
  - prefix: "--input_file"
    valueFrom: $(inputs.vcf)
  - "--check_existing"
  - prefix: "-o"
    valueFrom: stdout
  - prefix: "--assembly"
    valueFrom: $(inputs.assembly)
stdout: $(inputs.vcf.nameroot)_consequence.txt
