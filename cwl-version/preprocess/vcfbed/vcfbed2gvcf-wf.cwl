cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: vcfbed2homref
  - class: ResourceRequirement
    coresMin: 1
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
  vcfbed2gvcf:
    run: vcfbed2gvcf.cwl
    scatter: [vcf, bed, outname]
    scatterMethod: dotproduct
    in:
      vcf: get-vcfbed/vcfs
      bed: get-vcfbed/beds
      ref: ref
      outname: get-vcfbed/outnames
    out: [result]
