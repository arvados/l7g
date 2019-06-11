$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Creates gVCF with a given VCF, BED and reference FASTA
requirements:
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
  - class: ResourceRequirement
    coresMin: 1
baseCommand: bash
inputs:
  script:
    type: File
    label: Script to run vcfbed2homref, compress and index VCF
    default:
      class: File
      location: src/convert-vcf-bed-to-gvcf
    inputBinding:
      position: 1
  vcf:
    type: File
    label: VCF to be converted to gVCF
    inputBinding:
      position: 2
    secondaryFiles:
      - .tbi
  bed:
    type: File
    label: BED representing called region of VCF
    inputBinding:
      position: 3
  ref:
    type: File
    label: Compressed FASTA reference
    inputBinding:
      position: 4
  outname:
    type: string
    label: String to maintain VCF naming convention for gVCF
    inputBinding:
      position: 5
outputs:
  result:
    type: File
    label: Compressed gVCF and index file
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles:
      - .tbi
