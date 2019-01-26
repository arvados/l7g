$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a cgf for each FastJ file
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ScatterFeatureRequirement
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  refdirectory:
    type: Directory
    label: Input directory of FastJs
  cglf:
    type: Directory
    label: Tile library directory

outputs:
  out1:
    type: File[]
    label: Output cgfs
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
