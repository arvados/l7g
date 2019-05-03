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
  gvcfdir:
    type: Directory
    label: gVCF directory
  checknum:
    type: int
    label: Number of samples to check
  chroms:
    type: string[]
    label: Chromosomes to analyze
  tileassembly:
    type: File
    label: Reference tile assembly file
    secondaryFiles: [.fwi, .gzi]
  reffa:
    type: File
    label: Reference FASTA file
    secondaryFiles: [.fai, .gzi]

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

  createcgf-wf:
    run: ../cgf3/createcgf-wf.cwl
    in:
      fjdir: fjdir
      lib: merge-tilelib/mergedlib
    out: [cgfs]

  handle-cgfs:
    run: expressiontool/array-to-dir.cwl
    in:
      arr: createcgf-wf/cgfs
      dirname:
        valueFrom: "cgf"
    out: [dir]

  check-cgf-gvcf-wf:
    run: ../checks/check-cgf/gvcf/check-cgf-gvcf-wf.cwl
    in:
      cgfdir: handle-cgfs/dir
      sglfdir: merge-tilelib/mergedlib
      gvcfdir: gvcfdir
      checknum: checknum
      chroms: chroms
      tileassembly: tileassembly
      reffa: reffa
    out: [logs]

  createnpy-wf:
    run: ../npy/createnpy-wf.cwl
    in:
      waitsignal: check-cgf-gvcf-wf/logs
      cgfdir: handle-cgfs/dir
    out: [consolnpydir, names]
