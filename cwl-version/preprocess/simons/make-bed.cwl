$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: CommandLineTool
label: Make BED from VCF for regions passing a specified QUAL and GQ cutoff
requirements:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    coresMin: 2
    ramMin: 22000
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
  bed:
    type: stdout
    label: BED for regions that pass cutoff
baseCommand: [bcftools, view]
arguments:
  - prefix: "-e"
    valueFrom: "QUAL<$(inputs.qualcutoff) | QUAL='.' | FORMAT/GQ<$(inputs.gqcutoff)"
  - $(inputs.vcf)
  - shellQuote: false
    valueFrom: "|"
  - "convert2bed"
  - prefix: "-i"
    valueFrom: "vcf"
  - "-d"
  - shellQuote: false
    valueFrom: "|"
  - "cut"
  - "-f1-3"
  - shellQuote: false
    valueFrom: "|"
  - "bedtools"
  - "merge"
  - prefix: "-i"
    valueFrom: "-"
stdout: $(inputs.sample).bed
