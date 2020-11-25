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
  vcfsdir:
    type: Directory
    label: Input directory of VCFs
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
    outputSource: gvcf2fasta-wf/fas

steps:
  getfiles:
    run: getfiles.cwl
    in:
      dir: vcfsdir
    out: [vcfs, samples]
  gvcf2fasta-wf:
    run: gvcf2fasta-wf.cwl
    scatter: [sampleid, vcf]
    scatterMethod: dotproduct
    in:
      sampleid: getfiles/samples
      vcf: getfiles/vcfs
      genomebed: genomebed
      ref: ref
    out: [fas]
