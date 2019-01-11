$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a FASTJ file per path for each named GVCF
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  refdirectory:
    type: Directory
    label: Input Directory
  datafilenames:
    type: File
    label: gVCFs to create FastJ for
  bashscript:
    type: File
    label: Master script to create a FastJ for each gVCF
  ref:
    type: string
    label: Reference genome
  reffa:
    type: File
    label: Reference genome in FASTA format
  afn:
    type: File
    label: Compressed assembly fixed width file
  aidx:
    type: File
    label: Assembly index file
  refM:
    type: string
    label: Mitochondrial reference genome
  reffaM:
    type: File
    label: Reference mitochondrial genome in FASTA format
  afnM:
    type: File
    label: Compressed mitochondrial assembly fixed width file
  aidxM:
    type: File
    label: Mitochondrial assembly index file
  seqidM:
    type: string
    label: Mitochondrial naming scheme for gVCF
  tagdir:
    type: File
    label: Compressed tagset in FASTA format
  l7g:
    type: File
    label: Lightning application for parsing and searching assembly files
  pasta:
    type: File
    label: Tool for streaming and converting variant call formats
  refstream:
    type: File
    label: Tool to stream from FASTA file
  tile_assembly:
    type: File
    label: Tool to extract information from the tile assembly files

outputs:
  out1:
    type: Directory[]
    outputSource: step2/out1
  out2:
    type:
      type: array
      items:
        type: array
        items: File
    outputSource: step2/out2

steps:
  step1:
    run: getdirs_testset.cwl
    in:
      datafilenames: datafilenames
      refdirectory: refdirectory
    out: [out1,out2]

  step2:
    scatter: [gffDir,gffPrefix]
    scatterMethod: dotproduct
    in:
      bashscript: bashscript
      gffDir: step1/out1
      gffPrefix: step1/out2
      ref: ref
      reffa: reffa
      afn: afn
      aidx: aidx
      refM: refM
      reffaM: reffaM
      afnM: afnM
      aidxM: aidxM
      seqidM: seqidM
      tagdir: tagdir
      l7g: l7g
      pasta: pasta
      refstream: refstream
      tile_assembly: tile_assembly
    run: convertgvcf.cwl
    out: [out1,out2]
