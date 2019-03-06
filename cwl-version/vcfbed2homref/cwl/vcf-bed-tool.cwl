$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: vcfbed2homref0.1.3
  - class: ResourceRequirement
    coresMin: 1
  - class: InlineJavascriptRequirement
  - class: InitialWorkDirRequirement
    listing:
     - entry: $(inputs.refFaFn)
       writable: True
     # - entry: $(inputs.refFaFn.secondaryFiles[0])
       # writable: True
hints:
  - class: arv:RuntimeConstraints

baseCommand: bash

inputs:

  script:
    type: File
    inputBinding:
      position: 1

  gvcfFn:
    type: File
    inputBinding:
      position: 2
    secondaryFiles:
      - .tbi

  bedFn:
    type: File
    inputBinding:
      position: 3

  refFaFn:
    type: File
    inputBinding:
      position: 4
      valueFrom: $(self.basename)
    # secondaryFiles:
      # - .gzi
      # - .fai

  outName:
    type: string
    inputBinding:
      position: 5

outputs:
  result:
    type: Directory
    outputBinding:
      glob: "."
