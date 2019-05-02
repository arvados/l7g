cwlVersion: v1.0
class: CommandLineTool
label: Validate the conversion of the gVCF to cgf
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 1
    ramMin: 65000
inputs:
  script:
    type: File
    label: Master script to run validation
    default:
      class: File
      location: src/verify-conversion-batch-gvcf-cgf_skip-empty-and-zero-tilepaths.sh
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
  chrom:
    type: string
    label: Chromosome to analyze
  tileassembly:
    type: File
    label: Reference tile assembly file
    secondaryFiles: [.fwi, .gzi]
  reffa:
    type: File
    label: Reference FASTA file
    secondaryFiles: [.fai, .gzi]
outputs:
  log:
    type: File
    label: Validation logs
    outputBinding:
      glob: "*output.log"
baseCommand: bash
arguments:
  - $(inputs.script)
  - $(inputs.cgfdir)
  - $(inputs.sglfdir)
  - $(inputs.gvcfdir)
  - $(inputs.checknum)
  - $(inputs.chrom)
  - $(inputs.tileassembly)
  - $(inputs.reffa)
