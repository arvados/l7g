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
    label: Master script to merge tile libraries
    inputBinding:
      position: 1
  srcdir:
    type: Directory
    label: Directory of existing SGLF files
    inputBinding:
      position: 2
  nppdir:
    type: Directory
    label: Directory of new SGLF files for addition to tile library
    inputBinding:
      position: 3
  nthreads:
    type: string
    label: Number of threads to use
    inputBinding:
      position: 4
  mergetilelib:
    type: File
    label: Code that merges SGLF libraries
    inputBinding:
      position: 5

outputs:
  out1:
    type: Directory
    label: Directory of merged compressed tile library
    outputBinding:
      glob: "*merge*"
