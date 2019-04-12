cwlVersion: v1.0
class: Workflow
label: Clean GFFs to remove header lines
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ScatterFeatureRequirement: {}
inputs:
  gffdir:
    type: Directory
    label: Directory of input GFFs

outputs:
  cleangffs:
    type: File[]
    label: Clean GFFs
    outputSource: clean-gff-header/cleangff

steps:
  getfiles:
    run: getfiles.cwl
    in:
      gffdir: gffdir
    out: [gffs]

  clean-gff-header:
    run: clean-gff-header.cwl
    scatter: gff
    in:
      gff: getfiles/gffs
    out: [cleangff]
