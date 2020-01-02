$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Fix VCF by processing chrM
requirements:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
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
outputs:
  fixedvcf:
    type: File
    label: Fixed VCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [rtg, vcffilter]
arguments:
  - prefix: "-i"
    valueFrom: $(inputs.vcf)
  - prefix: "-o"
    valueFrom: $(inputs.vcf.basename)
  - prefix: "--javascript"
    valueFrom: $(inputs.filterjs)
