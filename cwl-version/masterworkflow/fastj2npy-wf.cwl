$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}

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
  srclib:
    type: Directory?
    label: Existing tile library directory

outputs:
  lib:
    type: Directory
    label: Tile library directory
    outputSource: merge-tilelib/mergedlib
  cgfdir:
    type: Directory
    label: Output cgfs
    outputSource: handle-cgfs/dir
  consolnpydir:
    type: Directory
    label: Output consolidated NumPy arrays
    outputSource: createnpy-wf/consolnpydir
  names:
    type: File
    label: File listing sample names
    outputSource: createnpy-wf/names

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
    run: jstools/nestedarray-to-dir.cwl
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

  createcgf-wf:
    run: ../cgf3/createcgf-wf.cwl
    in:
      fjdir: fjdir
      lib: merge-tilelib/mergedlib
    out: [cgfs]

  handle-cgfs:
    run: jstools/array-to-dir.cwl
    in:
      arr: createcgf-wf/cgfs
      dirname:
        valueFrom: "cgf"
    out: [dir]

  createnpy-wf:
    run: ../npy/createnpy-wf.cwl
    in:
      cgfdir: handle-cgfs/dir
    out: [consolnpydir, names]
