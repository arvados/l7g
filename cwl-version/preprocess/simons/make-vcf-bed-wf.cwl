cwlVersion: v1.0
class: Workflow
label: Scatter to filter VCF and make BED region
requirements:
  ScatterFeatureRequirement: {}
inputs:
  variantsvcfdir:
    type: Directory
    label: Input variants only VCF directory
  fullvcfdir:
    type: Directory
    label: Input full VCF directory
  qualcutoff:
    type: int
    label: Filtering QUAL cutoff
  gqcutoff:
    type: int
    label: Filtering GQ cutoff

outputs:
  filteredvcfs:
    type: File[]
    label: Output VCFs
    outputSource: filter-vcf/filteredvcf
  beds:
    type: File[]
    label: Output BEDs
    outputSource: make-bed/bed

steps:
  getvariantsvcfs:
    run: getfiles.cwl
    in:
      dir: variantsvcfdir
    out: [vcfs, samples]
  getfullvcfs:
    run: getfiles.cwl
    in:
      dir: fullvcfdir
    out: [vcfs, samples]
  filter-vcf:
    run: filter-vcf.cwl
    scatter: [vcf, sample]
    scatterMethod: dotproduct
    in:
      vcf: getvariantsvcfs/vcfs
      sample: getvariantsvcfs/samples
      qualcutoff: qualcutoff
      gqcutoff: gqcutoff
    out: [filteredvcf]
  make-bed:
    run: make-bed.cwl
    scatter: [vcf, sample]
    scatterMethod: dotproduct
    in:
      vcf: getfullvcfs/vcfs
      sample: getfullvcfs/samples
      qualcutoff: qualcutoff
      gqcutoff: gqcutoff
    out: [bed]
