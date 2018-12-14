$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a FastJ file for each gVCF by path
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ResourceRequirement
    coresMin: 1
    coresMax: 1
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  refdirectory:
    type: Directory
    label: Input directory
  bashscript:
    type: File
    label: Master script to convert FastJ to gVCFs
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
    label: Lightning application
  pasta:
    type: File
    label: Tool for streaming and converting variant call formats
  refstream:
    type: File
    label: Wrapper around 'samtools faidx' to get a stream out of a FASTA file
  tile_assembly:
    type: File
    label: Tool to extract information from the tile assembly files

outputs:
  out1:
    type: Directory[]
    outputSource: step2/out1
    label: Location of gVCFs and an array of directories
  out2:
    type:
      type: array
      items:
        type: array
        items: File
    outputSource: step2/out2
    label: Output FastJ Files

steps:
  step1:
    run: getdirs.cwl
    in:
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
