$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Convert FastJs to npy arrays
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}
hints:
 cwltool:LoadListingRequirement:
   loadListing: shallow_listing

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
    secondaryFiles: [^.fwi, .gzi]
  ref:
    type: string
    label: Reference genome
  reffa:
    type: File
    label: Reference FASTA file

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

  check-cgf-gvcf-wf:
    run: ../checks/check-cgf/gvcf/check-cgf-gvcf-wf.cwl
    in:
      cgfdir: handle-cgfs/dir
      sglfdir: merge-tilelib/mergedlib
      gvcfdir: gvcfdir
      checknum: checknum
      chroms: chroms
      tileassembly: tileassembly
      ref: ref
      reffa: reffa
    out: [gvcfhashes]

  createnpy-wf:
    run: ../npy/createnpy-wf.cwl
    in:
      waitsignal: check-cgf-gvcf-wf/gvcfhashes
      cgfdir: handle-cgfs/dir
    out: [consolnpydir, names]
