$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.1
class: Workflow
label: Creates a cgf for each FastJ file
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ScatterFeatureRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: no_listing

inputs:
  waitsignal:
    type: Any?
    label: Wait signal to start workflow
  fjdir:
    type: Directory
    label: Input directory of FastJs
  lib:
    type: Directory
    label: Tile library directory
  skippaths:
    type: File
    label: Paths to skip

outputs:
  cgfs:
    type: File[]
    label: Output cgfs
    outputSource: createcgf/cgf

steps:
  getdirs:
    run: getdirs.cwl
    in:
      fjdir: fjdir
    out: [fjdirs]

  createcgf:
    run: createcgf.cwl
    scatter: fjdir
    in:
      fjdir: getdirs/fjdirs
      lib: lib
      skippaths: skippaths
    out: [cgf]
