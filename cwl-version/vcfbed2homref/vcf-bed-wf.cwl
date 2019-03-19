cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: vcfbed2homref
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
  - class: ScatterFeatureRequirement
inputs:
  vcfsdir: Directory
  script: File
  refFaFn:
    type: File
outputs:
  result:
    type: Directory[]
    outputSource: vcftogvcftool/result
steps:
  getfiles:
    run: vcf-bed-scatter.cwl
    in: 
      vcfsdir: vcfsdir
    out: [vcfs, beds, out_files]
  vcftogvcftool:
    run: vcf-bed-tool.cwl
    scatter: [vcf, bed, out_file]
    scatterMethod: dotproduct
    in:
      script: script
      vcf: getfiles/vcfs
      bed: getfiles/beds
      ref: ref
      out_file: getfiles/out_files
    out: [result]