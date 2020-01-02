$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Fix VCF by changing haploid calls and processing chrM
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
  filterjs:
    type: File
    label: Javascript code for filtering
    default:
      class: File
      location: change_gt.js
outputs:
  fixedvcf:
    type: File
    label: Fixed VCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: zcat
arguments:
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "grep"
  - "-v"
  - "Locus GQX is less than 6 for hom deletion"
  - shellQuote: False
    valueFrom: "|"
  - "rtg"
  - "vcffilter"
  - prefix: "-i"
    valueFrom: "-"
  - prefix: "-o"
    valueFrom: $(inputs.vcf.basename)
  - prefix: "--javascript"
    valueFrom: $(inputs.filterjs)
