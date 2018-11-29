$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Merge all of the tile libraries (SGLFs)
doc: |
    Merge the compressed genome format files into a tile library (SGLF format)
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
    inputBinding:
      position: 1
    label: Bash script to merge cgf into Tile library
  srcdir:
    type: Directory
    inputBinding:
      position: 2
    label: Directory in keep where cgf files reside
  nppdir:
    type: Directory
    inputBinding:
      position: 3
  nthreads:
    type: string
    inputBinding:
      position: 4
    label: Number of threads to use
  mergetilelib:
    type: File
    inputBinding:
      position: 5
    label: Compiled C++ that reads in an SGLF line and stores the tilepath, tile library version, tilestep and tile span
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "*merge*"
