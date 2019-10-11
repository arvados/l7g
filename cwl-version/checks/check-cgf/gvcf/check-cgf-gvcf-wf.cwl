$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Workflow to validate the the gVCF to cgf conversion
requirements:
  ScatterFeatureRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: no_listing

inputs:
  cgfdir:
    type: Directory
    label: Compact genome format directory
  sglfdir:
    type: Directory
    label: Tile library directory
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
  gvcfhashes:
    type: File[]
    label: Hashes of gvcf pasta stream by path
    outputSource: check-cgf-gvcf/gvcfhash

steps:
  check-cgf-gvcf:
    run: check-cgf-gvcf.cwl
    scatter: chrom
    in:
      cgfdir: cgfdir
      sglfdir: sglfdir
      gvcfdir: gvcfdir
      checknum: checknum
      chrom: chroms
      tileassembly: tileassembly
      ref: ref
      reffa: reffa
    out: [gvcfhash]
