$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ScatterFeatureRequirement
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory: Directory
  cglf: Directory

outputs:
  out1:
    type: File[]
    outputSource: step2/out1
steps:
  step1:
    run: getdirs.cwl
    in:
      refdirectory: refdirectory
    out: [out1]
  step2:
    scatter: fjdir
    in:
      fjdir: step1/out1
      cglf: cglf
    run: createcgf.cwl
    out: [out1]
