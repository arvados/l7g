cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: vcfbed2homref0.1.3
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
    #secondaryFiles:
      #- .fai
      #- .gzi

outputs:
#   outNames:
#     type: File[]
#     outputSource: getfiles/gvcfFns
  result:
    type: Directory[]
    outputSource: vcftogvcftool/result
     

steps:
  getfiles:
    run: vcf-bed-scatter.cwl
    in: 
      vcfsdir: vcfsdir
    out: [gvcfFns, bedFns, outNames]
  vcftogvcftool:
    run: vcf-bed-tool.cwl
    scatter: [ gvcfFn, bedFn, outName ]
    scatterMethod: dotproduct
    in:
      script: script
      gvcfFn: getfiles/gvcfFns
      bedFn: getfiles/bedFns
      refFaFn: refFaFn
      outName: getfiles/outNames
    out: [result]
