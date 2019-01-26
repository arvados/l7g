$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Merges new tile library into existing tile library
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 16
    ramMin: 120000
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
    default:
      class: File
      location: ../src/merge-tilelibCWL.sh
    inputBinding:
      position: 1
  srcdir:
    type: Directory
    label: Existing tile library directory
    inputBinding:
      position: 2
  nppdir:
    type: Directory
    label: Directory of new tile library additions
    inputBinding:
      position: 3
  nthreads:
    type: string
    label: Number of threads to use
    inputBinding:
      position: 4
  mergetilelib:
    type: string
    label: Code that merges SGLF libraries
    default: "/usr/local/bin/merge-sglf"
    inputBinding:
      position: 5

outputs:
  out1:
    type: Directory
    label: Directory of merged tile library
    outputBinding:
      glob: "*merge*"
