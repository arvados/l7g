cwlVersion: v1.0
class: Workflow
label: Workflow to validate the the gVCF to cgf conversion
requirements:
  ScatterFeatureRequirement: {}

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
    secondaryFiles: [.fwi, .gzi]
  refFaFn:
    type: File
    label: Reference FASTA file
    secondaryFiles: [.fai, .gzi]
outputs:
  logs:
    type: File[]
    label: Validation logs
    outputSource: check-cgf-gvcf/log

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
      refFaFn: refFaFn
    out: [log]
