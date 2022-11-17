$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.1
class: Workflow
label: Scatter to Convert gVCF to FASTA
requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  arv:IntermediateOutput:
    outputTTL: 604800

inputs:
  sampleids:
    type: string[]
    label: Sample IDs
  vcftars:
    type: File[]
    label: Input VCF tars
  gqcutoff:
    type: int
    label: GQ (Genotype Quality) cutoff for filtering
  genomebed:
    type: File
    label: Whole genome BED
  ref:
    type: File
    label: Reference FASTA

outputs:
  fas:
    type:
      type: array
      items:
        type: array
        items: File
    label: Output pairs of FASTAs
    outputSource: gvcf2fasta_splitvcf-wf/fas

steps:
  gvcf2fasta_splitvcf-wf:
    run: gvcf2fasta_splitvcf-wf.cwl
    scatter: [sampleid, vcftar]
    scatterMethod: dotproduct
    in:
      sampleid: sampleids
      vcftar: vcftars
      gqcutoff: gqcutoff
      genomebed: genomebed
      ref: ref
    out: [fas]
