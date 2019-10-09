$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Convert GFFs to FastJ
requirements:
  ScatterFeatureRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

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

outputs:
  fjdirs:
    type: Directory[]
    label: Output FastJ directories
    outputSource: gff2fastj/fjdir

steps:
  getfiles:
    run: getfiles.cwl
    in:
      gffdir: gffdir
    out: [gffs]

  gff2fastj:
    run: gff2fastj.cwl
    scatter: gff
    in:
      gff: getfiles/gffs
      ref: ref
      reffa: reffa
      afn: afn
      tagset: tagset
      chroms: chroms
    out: [fjdir]
