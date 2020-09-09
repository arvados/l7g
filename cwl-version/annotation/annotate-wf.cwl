cwlVersion: v1.1
class: Workflow
label: Scatter paths to annotate tile library
requirements:
  ScatterFeatureRequirement: {}
inputs:
  ref:
    type: File
    label: Reference genome FASTA
  tilelib:
    type: Directory
    label: Input tile library
  assembly:
    type: File
    label: Compressed assembly fixed width file
    secondaryFiles: [^.fwi, .gzi]

outputs:
  annotations:
    type: File[]
    label: HGVS annotations in csv format
    outputSource: annotate/annotation

steps:
  getpaths:
    run: getpaths.cwl
    in:
      tilelib: tilelib
    out: [pathstrs]
  annotate:
    run: annotate.cwl
    scatter: pathstr
    in:
      pathstr: getpaths/pathstrs
      ref: ref
      tilelib: tilelib
      assembly: assembly
    out: [annotation]
