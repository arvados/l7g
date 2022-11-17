cwlVersion: v1.1
class: Workflow
requirements:
  ScatterFeatureRequirement: {}
inputs:
  sample: string
  chrs:
    type: string[]
  refsdir: Directory
  mapsdir: Directory
  vcf:
    type: File
    secondaryFiles: [.tbi]

outputs:
  rawimputedvcf:
    type: File
    outputSource: bcftools-concat/vcf

steps:
  match-ref-map-chr:
    run: match-ref-map-chr.cwl
    in:
      chrs: chrs
      refsdir: refsdir
      mapsdir: mapsdir
    out: [refs, maps]
  beagle:
    scatter: [chr, ref, map]
    scatterMethod: dotproduct
    run: beagle.cwl
    in:
      sample: sample
      chr: chrs
      ref: match-ref-map-chr/refs
      map: match-ref-map-chr/maps
      vcf: vcf
    out: [rawimputedvcf]
  bcftools-concat:
    run: bcftools-concat.cwl
    in:
      sample: sample
      vcfs: beagle/rawimputedvcf
    out: [vcf]
