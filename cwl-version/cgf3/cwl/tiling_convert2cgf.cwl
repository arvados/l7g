$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a cgf for each FASTJ file
doc: |
    Second intermediate step that takes in FASTJ files and creates compact genome representations of them for the Tile library.
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
  refdirectory
    type: Directory
    label: Keep input directory
  bashscript
    type: File
    label: Run convertcgfCWL-empty-problem-tilepaths3.sh
  cgft
    type: File
    label: cgft files
  fjt
    type: File
    label: JastJ Files
  cglf
    type: Directory
    label: Merged Directory

outputs:
  out1:
    type: File[]
    outputSource: step2/out1
    label: Outputs cgf from FASTJ
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
