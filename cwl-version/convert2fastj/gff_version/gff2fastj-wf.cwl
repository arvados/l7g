$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
  ScatterFeatureRequirement: {}
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  gffdir: Directory
  ref: string 
  reffa: File
  afn: File
  aidx: File
  refM: string
  reffaM: File
  afnM: File
  aidxM: File
  seqidM: string
  tagset: File

outputs:
  fjdirs:
    type: Directory[]
    outputSource: gff2fastj/fjdir

steps:
  getfiles:
    run: getfiles.cwl
    in:
      gffdir: gffdir
    out: [gffs]

  clean-gff-header:
    run: clean-gff-header.cwl
    scatter: gff
    in:
      gff: getfiles/gffs
    out: [cleangff]

  gff2fastj:
    run: gff2fastj.cwl
    scatter: gff
    in:
      gff: clean-gff-header/cleangff
      ref: ref
      reffa: reffa
      afn: afn
      aidx: aidx
      refM: refM
      reffaM: reffaM
      afnM: afnM
      aidxM: aidxM
      seqidM: seqidM
      tagset: tagset
    out: [fjdir]
