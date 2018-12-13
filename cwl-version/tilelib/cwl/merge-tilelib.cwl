$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Merge all of the tile libraries (SGLFs)
requirements:
  - class: DockerRequirement
    dockerPull: javatools-parallel
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    ramMin: 120000
    coresMin: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
baseCommand: bash
inputs:
  bashscriptmain:
    type: File
    label: Script to merge entries into tile library
    inputBinding:
      position: 1
  srcdir:
    type: Directory
    label: Directory of SGLF files
    inputBinding:
      position: 2
  nppdir:
    type: Directory
    label: Directory for new additions
    inputBinding:
      position: 3
  nthreads:
    type: string
    label: Number of threads to use
    inputBinding:
      position: 4
  mergetilelib:
    type: File
    label: Tool that takes an SGLF line and stores the tile path, tile library version, tile step and tile span
    inputBinding:
      position: 5

outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "*merge*"
