cwlVersion: v1.0
class: CommandLineTool
label: RTG vcfeval to compare VCFs
hints:
  DockerRequirement:
    dockerPull: vcfutil
inputs:
  baselinevcfgz:
    type: File
    label: Baseline VCF
    secondaryFiles: [.tbi]
  callsvcfgz:
    type: File
    label: Calls VCF
    secondaryFiles: [.tbi]
  sdf:
    type: Directory
    label: RTG reference directory
outputs:
  summary:
    type: File
    label: Summary file
    outputBinding:
      glob: "eval/summary.txt"
baseCommand: [rtg, vcfeval]
arguments:
  - prefix: "-b"
    valueFrom: $(inputs.baselinevcfgz)
  - prefix: "-c"
    valueFrom: $(inputs.callsvcfgz)
  - prefix: "-t"
    valueFrom: $(inputs.sdf)
  - prefix: "-o"
    valueFrom: "eval"
