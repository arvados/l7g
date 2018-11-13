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
  refdirectory: Directory
  datafilenames: File
  bashscript: File
  ref: string 
  reffa: File
  afn: File
  aidx: File
  refM: string
  reffaM: File
  afnM: File
  aidxM: File
  seqidM: string
  tagdir: File
  l7g: File
  pasta: File
  refstream: File
  tile_assembly: File

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
