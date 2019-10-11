$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Convert GFFs to npy arrays
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  SubworkflowFeatureRequirement: {}
  StepInputExpressionRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: no_listing

inputs:
  gffdir:
    type: Directory
    label: Input GFF directory
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
  gff2fastj-wf:
    run: ../convert2fastj/gff_version/gff2fastj-wf.cwl
    in:
      gffdir: gffdir
      ref: ref
      reffa: reffa
      afn: afn
      tagset: tagset
      chroms: chroms
    out: [fjdirs]

  handle-fjdirs:
    run: expressiontool/array-to-dir.cwl
    in:
      arr: gff2fastj-wf/fjdirs
      dirname:
        valueFrom: "fjdir"
    out: [dir]

  fastj2cgf-wf:
    run: fastj2cgf-wf.cwl
    in:
      fjdir: handle-fjdirs/dir
      tagset: tagset
      pathmin: pathmin
      pathmax: pathmax
      nchunks: nchunks
      sglfthreshold: sglfthreshold
      srclib: srclib
    out: [lib, sglfsize, skippaths, cgfdir]

  createnpy-wf:
    run: ../npy/createnpy-wf.cwl
    in:
      cgfdir: fastj2cgf-wf/cgfdir
    out: [consolnpydir, names]
