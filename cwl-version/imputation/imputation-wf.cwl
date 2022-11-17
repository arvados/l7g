cwlVersion: v1.1
class: Workflow
requirements:
  SubworkflowFeatureRequirement: {}
inputs:
  sample: string
  chrs:
    type: string[]
  refsdir: Directory
  mapsdir: Directory
  vcf:
    type: File
    secondaryFiles: [.tbi]
  nocallbed: File
  panelnocallbed: File
  panelcallbed: File
  genomebed: File

outputs:
  phasedimputedvcf:
    type: File
    outputSource: merge-phased-imputed-wf/phasedimputedvcf
  phasedimputednocallbed:
    type: File
    outputSource: merge-phased-imputed-wf/phasedimputednocallbed

steps:
  rtg-vcffilter:
    run: rtg-vcffilter.cwl
    in:
      sample: sample
      vcf: vcf
      excludebed: nocallbed
    out: [filteredvcf]
  scatter-beagle-wf:
    run: scatter-beagle-wf.cwl
    in:
      sample: sample
      chrs: chrs
      refsdir: refsdir
      mapsdir: mapsdir
      vcf: rtg-vcffilter/filteredvcf
    out: [rawimputedvcf]
  merge-phased-imputed-wf:
    run: merge-phased-imputed-wf.cwl
    in:
      sample: sample
      vcf: rtg-vcffilter/filteredvcf
      nocallbed: nocallbed
      rawimputedvcf: scatter-beagle-wf/rawimputedvcf
      panelnocallbed: panelnocallbed
      panelcallbed: panelcallbed
      genomebed: genomebed
    out: [phasedimputedvcf, phasedimputednocallbed]
