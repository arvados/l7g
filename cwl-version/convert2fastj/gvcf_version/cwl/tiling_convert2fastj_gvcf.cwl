$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a FASTJ file for each GVCF
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
  refdirectory
    type: Directory
  bashscript
    type: File
  ref
    type: string
  reffa
    type: File
  afn
    type: File
  aidx
    type: File
  refM
    type: string
  reffaM
    type: File
  afnM
    type: File
  aidxM
    type: File
  seqidM
    type: string
  tagdir
    type: File
  l7g
    type: File
  pasta
    type: File
  refstream
    type: File
  tile_assembly
    type: File

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
