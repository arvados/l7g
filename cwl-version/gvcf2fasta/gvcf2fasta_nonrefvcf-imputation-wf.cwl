cwlVersion: v1.1
class: Workflow
label: Impute gVCF and convert to FASTA for gVCF with NON_REF
requirements:
  ScatterFeatureRequirement: {}
  SubworkflowFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil

inputs:
  sampleid:
    type: string
    label: Sample ID
  vcf:
    type: File
    label: Input gVCF
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering
  genomebed:
    type: File
    label: Whole genome BED
  ref:
    type: File
    label: Reference FASTA
  haplotypes:
    type: int[]
    label: Haplotypes of sample
    default: [1, 2]
  chrs:
    type: string[]
  refsdir: Directory
  mapsdir: Directory
  panelnocallbed: File

outputs:
  fas:
    type: File[]
    label: Output pair of FASTAs
    outputSource: bcftools-consensus/fa

steps:
  fixvcf-get_bed_varonlyvcf:
    run: fixvcf-get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      vcf: vcf
      gqcutoff: gqcutoff
      genomebed: genomebed
    out: [nocallbed, varonlyvcf]

  imputation-wf:
    run: ../imputation/imputation-wf.cwl
    in:
      sample: sampleid
      chrs: chrs
      refsdir: refsdir
      mapsdir: mapsdir
      vcf: fixvcf-get_bed_varonlyvcf/varonlyvcf
      nocallbed: fixvcf-get_bed_varonlyvcf/nocallbed
      panelnocallbed: panelnocallbed
    out: [phasedimputedvcf, phasedimputednocallbed]

  bcftools-consensus:
    run: bcftools-consensus.cwl
    scatter: haplotype
    in:
      sampleid: sampleid
      vcf: imputation-wf/phasedimputedvcf
      ref: ref
      haplotype: haplotypes
      mask: imputation-wf/phasedimputednocallbed
    out: [fa]
