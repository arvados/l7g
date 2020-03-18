$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.1
class: Workflow
label: Convert FastJs to npy arrays for gVCF input
requirements:
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: no_listing

inputs:
  gvcfdir:
    type: Directory
    label: Input gVCF directory
  fjdir:
    type: Directory
    label: Directory of FastJ files
  ref:
    type: string
    label: Reference genome
  reffa:
    type: File
    label: Reference genome in FASTA format
  afn:
    type: File
    label: Compressed assembly fixed width file
    secondaryFiles: [^.fwi, .gzi]
  tagset:
    type: File
    label: Compressed tagset in FASTA format
  chroms:
    type: string[]
    label: Chromosomes to analyze
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
  checknum:
    type: int
    label: Number of samples to check
  checkchroms:
    type: string[]
    label: Chromosomes to validate

outputs:
  lib:
    type: Directory
    label: Tile library directory
    outputSource: fastj2cgf-wf/lib
  sglfsize:
    type: File
    label: Unzipped sglf size
    outputSource: fastj2cgf-wf/sglfsize
  skippaths:
    type: File
    label: Paths to skip
    outputSource: fastj2cgf-wf/skippaths
  cgfdir:
    type: Directory
    label: Output cgfs
    outputSource: fastj2cgf-wf/cgfdir
  consolnpydir:
    type: Directory
    label: Output consolidated NumPy arrays
    outputSource: createnpy-wf/consolnpydir
  names:
    type: File
    label: File listing sample names
    outputSource: createnpy-wf/names

steps:
  fastj2cgf-wf:
    run: fastj2cgf-wf.cwl
    in:
      fjdir: fjdir
      tagset: tagset
      pathmin: pathmin
      pathmax: pathmax
      nchunks: nchunks
      sglfthreshold: sglfthreshold
      srclib: srclib
    out: [lib, sglfsize, skippaths, cgfdir]

  check-cgf-gvcf-wf:
    run: ../checks/check-cgf/gvcf/check-cgf-gvcf-wf.cwl
    in:
      cgfdir: fastj2cgf-wf/cgfdir
      sglfdir: fastj2cgf-wf/lib
      gvcfdir: gvcfdir
      checknum: checknum
      chroms: checkchroms
      tileassembly: afn
      ref: ref
      reffa: reffa
    out: [gvcfhashes]

  createnpy-wf:
    run: ../npy/createnpy-wf.cwl
    in:
      waitsignal: check-cgf-gvcf-wf/gvcfhashes
      cgfdir: fastj2cgf-wf/cgfdir
    out: [consolnpydir, names]
