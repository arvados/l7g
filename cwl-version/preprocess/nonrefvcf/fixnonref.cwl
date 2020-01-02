$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Fix VCF by removing GT fields that point to <NON_REF> and processing chrM
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
    default:
      class: File
      location: ../chrmvcf/change_gt_chrM.js
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
  - prefix: "--keep-expr"
    valueFrom: "ALT.length == 1 || SAMPLES[0].GT.indexOf(String(ALT.length)) == -1"
  - prefix: "--javascript"
    valueFrom: $(inputs.filterjs)
