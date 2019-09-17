$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Validate the conversion of the gVCF to cgf
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4000
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
    secondaryFiles: [^.fwi, .gzi]
  ref:
    type: string
    label: Reference genome
  reffa:
    type: File
    label: Reference FASTA file
outputs:
  gvcfhash:
    type: File
    label: Hashes of gvcf pasta stream by path
    outputBinding:
      glob: "gvcf_hash_file"
baseCommand: bash
arguments:
  - $(inputs.script)
  - $(inputs.cgfdir)
  - $(inputs.sglfdir)
  - $(inputs.gvcfdir)
  - $(inputs.checknum)
  - $(inputs.chrom)
  - $(inputs.tileassembly)
  - $(inputs.ref)
  - $(inputs.reffa)
