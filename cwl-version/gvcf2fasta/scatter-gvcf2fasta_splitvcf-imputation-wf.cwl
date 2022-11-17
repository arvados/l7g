$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: Workflow
label: Scatter to impute gVCF and convert gVCF to FASTA
requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil

inputs:
  sampleids:
    type: string[]
    label: Sample IDs
  splitvcfdirs:
    type: Directory[]
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
    type:
      type: array
      items:
        type: array
        items: File
    label: Output pairs of FASTAs
    outputSource: gvcf2fasta_splitvcf-imputation-wf/fas

steps:
  gvcf2fasta_splitvcf-imputation-wf:
    run: gvcf2fasta_splitvcf-imputation-wf.cwl
    scatter: [sampleid, splitvcfdir]
    scatterMethod: dotproduct
    in:
      sampleid: sampleids
      splitvcfdir: splitvcfdirs
      gqcutoff: gqcutoff
      genomebed: genomebed
      ref: ref
      chrs: chrs
      refsdir: refsdir
      mapsdir: mapsdir
      panelnocallbed: panelnocallbed
      panelcallbed: panelcallbed
    out: [fas]
