$namespaces:
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Concatenate a set of VCFs split by chromosomes
requirements:
  ScatterFeatureRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  vcfdirs:
    type: Directory[]
    label: Input VCFs directories

outputs:
  vcfs:
    type: File[]
    label: Concatenated VCFs
    outputSource: zcatvcf/vcf
    secondaryFiles: [.tbi]

steps:
  zcatvcf:
    run: zcatvcf.cwl
    scatter: vcfdir
    in:
      vcfdir: vcfdirs
    out: [vcf]
