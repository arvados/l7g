$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Resolve duplicate/overlapping calls in the gVCFs in a named list
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ResourceRequirement
    coresMin: 1 
    coresMax: 1
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  refdirectory: Directory
  datafilenames: File
  bashscript: File
  cleanvcf: File

outputs:
  out1:
    type: Directory[]
    outputSource: step2/out1

steps:
  step1:
    run: getdirs_testset.cwl
    in: 
      datafilenames: datafilenames
      refdirectory: refdirectory 
    out: [out1,out2]

  step2:
    scatter: [gvcfDir,gvcfPrefix] 
    scatterMethod: dotproduct
    in: 
      bashscript: bashscript
      gvcfDir: step1/out1
      gvcfPrefix: step1/out2
      cleanvcf: cleanvcf
    run: cleangvcf.cwl
    out: [out1]
