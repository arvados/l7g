$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: CommandLineTool
label: Filters VCF by a specified QUAL and GQ cutoff
requirements:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
  ShellCommandRequirement: {}
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  vcf:
    type: File
    label: Input VCF file
  sample:
    type: string
    label: Sample name of VCF
  qualcutoff:
    type: int
    label: Filtering QUAL cutoff
  gqcutoff:
    type: int
    label: Filtering GQ cutoff
outputs:
  filteredvcf:
    type: File
    label: Filtered VCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bcftools, view]
arguments:
  - "-Oz"
  - prefix: "-o"
    valueFrom: $(inputs.sample).vcf.gz
  - prefix: "-e"
    valueFrom: "QUAL<$(inputs.qualcutoff) | QUAL='.' | FORMAT/GQ<$(inputs.gqcutoff)"
  - $(inputs.vcf)
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - prefix: "-p"
    valueFrom: "vcf"
  - $(inputs.sample).vcf.gz
