cwlVersion: v1.1
class: CommandLineTool
label: Run SnpEff on given VCF and use bcftools to annotate with dbSNP and gnomAD
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: snpeff4.3t
  ResourceRequirement:
    coresMin: 2
    ramMin: 20000
    tmpdirMin: 16000
inputs:
  vcf:
    type: File
    label: Input VCF
  sample:
    type: string
    label: Sample name
  snpeffdatadir:
    type: Directory
    label: Database directory for SnpEff
  genomeversion:
    type: string
    label: Genome version
  dbsnp:
    type: File
    label: dbSNP BCF
    secondaryFiles: [.csi]
  gnomad:
    type: File
    label: gnomAD BCF
    secondaryFiles: [.csi]
outputs:
  annotatedvcf:
    type: File
    label: Annotated VCF
    outputBinding:
      glob: "*_snpeff_dbsnp_gnomad.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [java]
arguments:
  - -Xmx$(runtime.ram)m
  - prefix: "-jar"
    valueFrom: "/snpEff/snpEff.jar"
  - prefix: "-dataDir"
    valueFrom: $(inputs.snpeffdatadir)
  - $(inputs.genomeversion)
  - $(inputs.vcf)
  - shellQuote: False
    valueFrom: "|"
  - "bgzip"
  - "-c"
  - shellQuote: False
    valueFrom: ">"
  - $(inputs.sample)_snpeff.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_snpeff.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "bcftools"
  - "annotate"
  - prefix: "--annotations"
    valueFrom: $(inputs.dbsnp)
  - prefix: "--columns"
    valueFrom: "=ID"
  - $(inputs.sample)_snpeff.vcf.gz
  - "-Oz"
  - prefix: "-o"
    valueFrom: $(inputs.sample)_snpeff_dbsnp.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_snpeff_dbsnp.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "bcftools"
  - "annotate"
  - prefix: "--annotations"
    valueFrom: $(inputs.gnomad)
  - prefix: "--columns"
    valueFrom: "INFO/AC,INFO/AN,INFO/AF,INFO/AF_afr,INFO/AF_amr,INFO/AF_asj,INFO/AF_eas,INFO/AF_fin,INFO/AF_nfe,INFO/AF_oth"
  - $(inputs.sample)_snpeff_dbsnp.vcf.gz
  - "-Oz"
  - prefix: "-o"
    valueFrom: $(inputs.sample)_snpeff_dbsnp_gnomad.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_snpeff_dbsnp_gnomad.vcf.gz
  - shellQuote: False
    valueFrom: "&&"
  - "rm"
  - $(inputs.sample)_snpeff.vcf.gz
  - $(inputs.sample)_snpeff.vcf.gz.tbi
  - $(inputs.sample)_snpeff_dbsnp.vcf.gz
  - $(inputs.sample)_snpeff_dbsnp.vcf.gz.tbi
