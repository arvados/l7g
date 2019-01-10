$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a cgf for each FastJ file
requirements:
  - class: DockerRequirement
    dockerPull: javatoolsparallel
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement

hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory:
    type: Directory
    label: Input directory of FastJs
  bashscript:
    type: File
    label: Master script to convert FastJs to cgfs
  cgft:
    type: File
    label: Tool to manipulate and inspect cgf files
  fjt:
    type: File
    label: Tool to manipulate FastJ files
  cglf:
    type: Directory
    label: Tile library directory

outputs:
  out1:
    type: File[]
    outputSource: step2/out1
    label: Output cgfs
steps:
  step1:
    run: getdirs.cwl
    in:
      refdirectory: refdirectory
    out: [out1]
  step2:
    scatter: fjdir
    scatterMethod: dotproduct
    in:
      fjdir: step1/out1
      bashscript: bashscript
      cgft: cgft
      fjt: fjt
      cglf: cglf
    run: createcgf.cwl
    out: [out1]
