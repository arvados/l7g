cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: fbh/vcfpreprocess
  - class: ResourceRequirement
    coresMin: 1
    ramMin: 13000
  - class: ScatterFeatureRequirement
inputs:
  vcfsdir: Directory
  ref:
    type: File
outputs:
  result:
    type: File[]
    outputSource: vcfbed2gvcf/result
steps:
  get-vcfbed:
    run: get-vcfbed.cwl
    in: 
      vcfsdir: vcfsdir
    out: [vcfs, beds, outnames]
  intersect-vcfbed:
    run: intersect-vcfbed.cwl
    scatter: [vcf, bed]
    scatterMethod: dotproduct
    in:
      vcf: get-vcfbed/vcfs
      bed: get-vcfbed/beds
    out: [sortedvcf, sortedbed]  
  vcfbed2gvcf:
    run: vcfbed2gvcf.cwl
    scatter: [vcf, bed, outname]
    scatterMethod: dotproduct
    in:
      vcf: intersect-vcfbed/sortedvcf
      bed: intersect-vcfbed/sortedbed
      ref: ref
      outname: get-vcfbed/outnames
    out: [result]
