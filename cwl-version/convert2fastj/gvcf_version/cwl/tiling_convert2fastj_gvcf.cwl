$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Creates a FastJ file for each gVCF by path
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  refdirectory:
    type: Directory
    label: Directory of input gVCFs
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

outputs:
  out1:
    type: Directory[]
    label: Directories of FastJs
    outputSource: step2/out1
  out2:
    type:
      type: array
      items:
        type: array
        items: File
    label: Indexed and zipped gVCFs
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
    run: convertgvcf.cwl
    out: [out1,out2]
