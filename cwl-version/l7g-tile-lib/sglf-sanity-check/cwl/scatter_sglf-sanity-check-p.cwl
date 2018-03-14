cwlVersion: v1.0
class: Workflow
requirements:
  ScatterFeatureRequirement: {}

inputs:
  script: File
  listFiles: File[]
  sglfDir: Directory
  outLog: string[]
  nCore: string

outputs:
  outfiles:
    type: Directory
    outputSource: gather/out

steps:
  batch_check:
    run: sglf-sanity-check-p.cwl
    scatter: listFns
    in:
      script: script
      listFile: listFiles
      outFileName: outLog
      nCore: nCore
      sglfDir: sglfDir
    out: [result]
  gather:
    run: gather_hupgp-gff-to-fastj.cwl
    in:
      indirs: convert/result
    out: [out]

