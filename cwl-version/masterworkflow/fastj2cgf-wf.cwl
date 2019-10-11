$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Convert FastJs to cgfs
requirements:
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: no_listing

inputs:
  fjdir:
    type: Directory
    label: Directory of FastJ files
  tagset:
    type: File
    label: Compressed tagset in FASTA format
  pathmin:
    type: string
    label: Starting path in the tile library
  pathmax:
    type: string
    label: Last/Maximum path in the tile library
  nchunks:
    type: string
    label: Number of chunks to scatter
  sglfthreshold:
    type: int
    label: Threshold for unzipped sglf size in MiB
  srclib:
    type: Directory?
    label: Existing tile library directory

outputs:
  lib:
    type: Directory
    label: Tile library directory
    outputSource: merge-tilelib/mergedlib
  sglfsize:
    type: File
    label: Unzipped sglf size
    outputSource: getsglfsize/sglfsize
  skippaths:
    type: File
    label: Paths to skip
    outputSource: getsglfsize/skippaths
  cgfdir:
    type: Directory
    label: Output cgfs
    outputSource: handle-cgfs/dir

steps:
  createsglf-wf:
    run: ../tilelib/createsglf-wf.cwl
    in:
      pathmin: pathmin
      pathmax: pathmax
      nchunks: nchunks
      fjdir: fjdir
      tagset: tagset
    out: [sglfs]

  handle-sglfs:
    run: expressiontool/nestedarray-to-dir.cwl
    in:
      nestedarr: createsglf-wf/sglfs
      dirname:
        valueFrom: "sglf"
    out: [dir]

  merge-tilelib:
    run: ../tilelib/merge-tilelib.cwl
    in:
      srclib: srclib
      newlib: handle-sglfs/dir
    out: [mergedlib]

  sglf-sanity-check:
    run: ../checks/check-sglf/sglf-sanity-check.cwl
    in:
      sglfdir: merge-tilelib/mergedlib
    out: [log]

  getsglfsize:
    run: ../checks/check-sglf/getsglfsize.cwl
    in:
      lib: merge-tilelib/mergedlib
      threshold: sglfthreshold
    out: [sglfsize, skippaths]

  createcgf-wf:
    run: ../cgf3/createcgf-wf.cwl
    in:
      waitsignal: sglf-sanity-check/log
      fjdir: fjdir
      lib: merge-tilelib/mergedlib
      skippaths: getsglfsize/skippaths
    out: [cgfs]

  handle-cgfs:
    run: expressiontool/array-to-dir.cwl
    in:
      arr: createcgf-wf/cgfs
      dirname:
        valueFrom: "cgf"
    out: [dir]
