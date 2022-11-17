cwlVersion: v1.2
class: CommandLineTool
hints:
  DockerRequirement:
    dockerPull: vcfutil
  ResourceRequirement:
    ramMin: 5000
inputs:
  baselinevcf:
    type: File
    secondaryFiles: [.tbi]
  callsvcf:
    type: File
    secondaryFiles: [.tbi]
  sdf:
    type: Directory
outputs:
  evaldir:
    type: Directory
    outputBinding:
      glob: "eval"
baseCommand: [rtg, vcfeval]
arguments:
  - prefix: "-b"
    valueFrom: $(inputs.baselinevcf)
  - prefix: "-c"
    valueFrom: $(inputs.callsvcf)
  - prefix: "-t"
    valueFrom: $(inputs.sdf)
  - prefix: "-o"
    valueFrom: "eval"
