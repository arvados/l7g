cwlVersion: v1.0
class: CommandLineTool
label: Remove unused annotations
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
inputs:
  vcfgz:
    type: File
    label: Input VCF
outputs:
  annotatedvcfgz:
    type: File
    label: Annotated VCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bcftools, annotate]
arguments:
  - prefix: "-x"
    valueFrom: "INFO/customer_score1,INFO/customer_score2,INFO/ADP,INFO/ADP,INFO/HET,INFO/HOM,INFO/NC,INFO/WT,FORMAT/AO,FORMAT/GL,FORMAT/QA,FORMAT/SDP,FORMAT/RD,FORMAT/AD,FORMAT/FREQ,FORMAT/PVAL,FORMAT/RBQ,FORMAT/ABQ,FORMAT/RDF,FORMAT/RDR,FORMAT/ADF,FORMAT/ADR"
  - $(inputs.vcfgz)
  - prefix: "-O"
    valueFrom: "z"
  - prefix: "-o"
    valueFrom: $(inputs.vcfgz.basename)
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - prefix: "-p"
    valueFrom: "vcf"
  - $(inputs.vcfgz.basename)
