$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: Workflow
label: Preprocess portable VCF
requirements:
  arv:RunInSingleContainer: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 12000

inputs:
  vcfgz:
    type: File
    label: Input VCF
  header:
    type: File
    label: Header file
  sdf:
    type: Directory
    label: RTG reference directory
  cleanvcf:
    type: File
    label: Code that cleans VCFs

outputs:
  processedvcfgz:
    type: File
    label: Processed VCF
    outputSource: bcftools-annotate/annotatedvcfgz
  summary:
    type: File
    label: Summary file
    outputSource: rtg-vcfeval/summary

steps:
  bcftools-reheader:
    run: bcftools-reheader.cwl
    in:
      header: header
      vcfgz: vcfgz
    out: [reheaderedvcfgz]

  sort-clean:
    run: sort-clean.cwl
    in:
      vcfgz: bcftools-reheader/reheaderedvcfgz
      cleanvcf: cleanvcf
    out: [cleanvcfgz]

  bcftools-annotate:
    run: bcftools-annotate.cwl
    in:
      vcfgz: sort-clean/cleanvcfgz
    out: [annotatedvcfgz]

  rtg-vcfeval:
    run: rtg-vcfeval.cwl
    in:
      baselinevcfgz: bcftools-annotate/annotatedvcfgz
      callsvcfgz: bcftools-annotate/annotatedvcfgz
      sdf: sdf
    out: [summary]
