$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Intersect VCF and BED
requirements:
  - class: ShellCommandRequirement
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
  - class: ResourceRequirement
    ramMin: 12000
inputs:
  vcf:
    type: File
    label: VCF to be intersected
  bed:
    type: File
    label: BED to intersect with VCF
outputs:
  intersectedvcf:
    type: File
    label: Intersected VCF with 100% alignment
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bedtools, intersect]
arguments:
  - "-header"
  - prefix: "-a"
    valueFrom: $(inputs.vcf)
  - prefix: "-b"
    valueFrom: $(inputs.bed)
  - prefix: "-f"
    valueFrom: "1"
  - shellQuote: false
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: false
    valueFrom: ">"
  - $(inputs.vcf.basename)
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - prefix: "-p"
    valueFrom: "vcf"
  - $(inputs.vcf.basename)
