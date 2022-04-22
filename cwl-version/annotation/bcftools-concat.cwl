$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: CommandLineTool
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: snpeff4.3t
  ResourceRequirement:
    coresMin: 2
    ramMin: 10000
  arv:RuntimeConstraints:
    keep_cache: 20000
    outputDirType: keep_output_dir
inputs:
  sample: string
  vcfs:
    type: File[]
    secondaryFiles: [.tbi]
outputs:
  vcf:
    type: File
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bcftools, concat]
arguments:
  - $(inputs.vcfs)
  - "-Oz"
  - prefix: "-o"
    valueFrom: $(inputs.sample)_snpeff_dbsnp_gnomad.vcf.gz
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_snpeff_dbsnp_gnomad.vcf.gz
