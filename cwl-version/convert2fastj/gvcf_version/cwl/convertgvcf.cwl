$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
baseCommand: bash
inputs:
  bashscript
    type: File
    label: Master bash script that controls converting FastJ to gVCFs
    inputBinding:
      position: 1
  gffDir:
    type: Directory
    label: Path with compressed gVCF files
    inputBinding:
      position: 2
  gffPrefix:
    type: string
    label: Prefix of all gVCF files
    inputBinding:
      position: 3
  ref
    type: string
    label: Reference genome
    inputBinding:
      position: 4
  reffa:
    type: File
    label: Reference genome in Fasta format
    inputBinding:
      position: 5
  afn:
    type: File
    label: Compressed assembly fixed width file
    inputBinding:
      position: 6
  aidx:
    type: File
    label: Assembly index file
    inputBinding:
      position: 7
  refM:
    type: string
    label: Mitochondrial reference genome
    inputBinding:
      position: 8
  reffaM:
    type: File
    label: Reference mitochondrial genome in FASTA format
    inputBinding:
      position: 9
  afnM:
    type: File
    label: Compressed mitochondrial assembly fixed width file
    inputBinding:
      position: 10
  aidxM:
    type: File
    label: Mitochondrial assembly index file
    inputBinding:
      position: 11
  seqidM:
    type: string
    label: Mitochondrial naming scheme for gVCF
    inputBinding:
      position: 12
  tagdir:
    type: File
    label: Compressed tagset in FASTA format
    inputBinding:
      position: 13
  l7g:
    type: File
    label: Lightning application
    inputBinding:
      position: 14
  pasta:
    type: File
    label: Tool for steaming and converting variant call formats
    inputBinding:
      position: 15
  refstream:
    type: File
    label: Wrapper around 'samtools faidx' to get a stream out of a FASTA file
    inputBinding:
      position: 16
  tile_assembly:
    type: File
    label: Tool to extract information from the tile assembly files
    inputBinding:
      position: 17
outputs:
  out1:
    type: Directory
    outputBinding:
      glob: "stage/*"
  out2:
    type: File[]
    outputBinding:
      glob: "indexed/*.gz*"
