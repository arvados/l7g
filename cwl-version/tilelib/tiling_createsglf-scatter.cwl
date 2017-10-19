$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    coresMin: 1 
    coresMax: 1 
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  pathmin: string
  pathmax: string
  bashscript: File
  fastj2cgflib: File
  datadir: File
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
       bashscript: bashcript
       tilepath: step1/out1
       fastj2cgflib: fastj2cgflib
       datadir: datadir
       verbose_target: verbose_target
       tagset: tagset:
     run: createsglfSingle.cwl
     out: [out1]
