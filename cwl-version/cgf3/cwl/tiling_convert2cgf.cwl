$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a cgf for each FastJ file
doc: |
    Takes in FastJ files and creates compact genome representations of them for the tile library.
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
    label: Bash script to convert FastJ to cgf using SGLF library
  cgft:
    type: File
    label: Compact genome format tool, for manipulating and inspecting cgf files
  fjt:
    type: File
    label: Tool to manipulate FastJ (text) files.
  cglf:
    type: Directory
    label: Tile library location

outputs:
  out1:
    type: File[]
    outputSource: step2/out1
    label: cgf created from FastJ
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
