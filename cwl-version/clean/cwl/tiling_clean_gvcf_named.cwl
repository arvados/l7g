$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Resolve duplicate/overlapping calls in the gVCFs in a named list
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  refdirectory: Directory
  datafilenames: File

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
      gvcfDir: step1/out1
      gvcfPrefix: step1/out2
    run: cleangvcf.cwl
    out: [out1]
