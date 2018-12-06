$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
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
    inputBinding:
      position: 1
  srcdir:
    type: Directory
    inputBinding:
      position: 2
  nppdir:
    type: Directory
    inputBinding:
      position: 3
  nthreads:
    type: string
    inputBinding:
      position: 4
  mergetilelib:
    type: [File,string]
    default: "/usr/local/bin/merge-sglf"
    inputBinding:
      position: 5
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "*merge*"
