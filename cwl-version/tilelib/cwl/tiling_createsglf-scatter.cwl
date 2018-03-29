$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow 
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
  - class: ScatterFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 16384 
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  pathmin: string
  pathmax: string
  bashscript: File
  fastj2cgflib: File 
  datadir: Directory 
  verbose_tagset: File
  tagset: File
outputs:
  out1: 
    type: File[]
    outputSource: step2/out1

steps:
  step1:
    run: getpaths.cwl
    in:
      pathmin: pathmin
      pathmax: pathmax
    out: [out1]
  
  step2:
     scatter: [tilepath]
     scatterMethod: dotproduct
     in: 
       bashscript: bashscript
       tilepath: step1/out1
       fastj2cgflib: fastj2cgflib
       datadir: datadir
       verbose_tagset: verbose_tagset
       tagset: tagset
     run: createsglfSingle.cwl
     out: [out1]
