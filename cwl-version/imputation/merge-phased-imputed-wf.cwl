cwlVersion: v1.1
class: Workflow
inputs:
  sample: string
  vcf:
    type: File
    secondaryFiles: [.tbi]
  nocallbed: File
  rawimputedvcf:
    type: File
    secondaryFiles: [.tbi]
  panelnocallbed: File
  panelcallbed: File
  genomebed: File

outputs:
  phasedimputedvcf:
    type: File
    outputSource: rtg-vcfmerge/phasedimputedvcf
  phasedimputednocallbed:
    type: File
    outputSource: bedtools-intersect_phasedimputednocallbed/intersectbed

steps:
  get-phasedvcf:
    run: get-phasedvcf.cwl
    in:
      sample: sample
      vcf: rawimputedvcf
    out: [phasedvcf]
  get-imputedvcf:
    run: get-imputedvcf.cwl
    in:
      sample: sample
      vcf: rawimputedvcf
    out: [imputedvcf]
  bedtools-intersect_phasedimputednocallbed:
    run: bedtools-intersect.cwl
    in:
      sample: sample
      a: nocallbed
      b: panelnocallbed
    out: [intersectbed]
  bedtools-intersect_imputationbed:
    run: bedtools-intersect.cwl
    in:
      sample: sample
      a: nocallbed
      b: panelcallbed
    out: [intersectbed]
  rtg-vcffilter-bedtools-intersect:
    run: rtg-vcffilter-bedtools-intersect.cwl
    in:
      sample: sample
      vcf: get-imputedvcf/imputedvcf
      bed: bedtools-intersect_imputationbed/intersectbed
    out: [filteredvcf]
  rtg-vcfmerge:
    run: rtg-vcfmerge.cwl
    in:
      sample: sample
      vcf: vcf
      phasedvcf: get-phasedvcf/phasedvcf
      imputedvcf: rtg-vcffilter-bedtools-intersect/filteredvcf
    out: [phasedimputedvcf]
