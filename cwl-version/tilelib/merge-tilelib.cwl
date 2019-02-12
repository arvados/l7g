$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Merges new tile library into existing tile library
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 16
    ramMin: 120000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
baseCommand: bash
inputs:
  bashscript:
    type: File
    label: Master script to merge tile libraries
    default:
      class: File
      location: src/merge-tilelibCWL.sh
    inputBinding:
      position: 1
  srclib:
    type: Directory?
    label: Existing tile library directory
    inputBinding:
      prefix: -s
      position: 2
  newlib:
    type: Directory
    label: New tile library directory to be added
    inputBinding:
      prefix: -n
      position: 3
  nthreads:
    type: string
    label: Number of threads to use
    default: "6"
    inputBinding:
      position: 4
  mergetilelib:
    type: string
    label: Code that merges SGLF libraries
    default: "/usr/local/bin/merge-sglf"
    inputBinding:
      position: 5
outputs:
  mergedlib:
    type: Directory
    label: Directory of merged tile library
    outputBinding:
      glob: "*merge*"
