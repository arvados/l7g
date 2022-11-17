cwlVersion: v1.1
class: Workflow
label: Impute gVCF and convert to FASTA for gVCF split by chromosome
requirements:
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil

inputs:
  sampleid:
    type: string
    label: Sample ID
  splitvcfdir:
    type: Directory
    label: Input directory of split gVCFs
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering
  genomebed:
    type: File
    label: Whole genome BED
  ref:
    type: File
    label: Reference FASTA
  chrs: string[]
  refsdir: Directory
  mapsdir: Directory
  panelnocallbed: File
  panelcallbed: File

outputs:
  fas:
    type: File[]
    label: Output pair of FASTAs
    outputSource: bcftools-consensus/fas

steps:
  concat-get_bed_varonlyvcf:
    run: concat-get_bed_varonlyvcf.cwl
    in:
      sampleid: sampleid
      splitvcfdir: splitvcfdir
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
      vcf: concat-get_bed_varonlyvcf/varonlyvcf
      nocallbed: concat-get_bed_varonlyvcf/nocallbed
      panelnocallbed: panelnocallbed
      panelcallbed: panelcallbed
      genomebed: genomebed
    out: [phasedimputedvcf, phasedimputednocallbed]

  append-sampleid:
    run: append-sampleid.cwl
    in:
      sampleid: sampleid
      suffix:
        valueFrom: "_phased_imputed"
    out: [appendedsampleid]

  bcftools-consensus:
    run: bcftools-consensus.cwl
    in:
      sampleid: append-sampleid/appendedsampleid
      vcf: imputation-wf/phasedimputedvcf
      ref: ref
      mask: imputation-wf/phasedimputednocallbed
    out: [fas]
