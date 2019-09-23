#!/usr/bin/env cwl-runner
$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
baseCommand: mergefastj

hints:
  DockerRequirement:
    dockerPull: jchen/golangmergefastj
  ResourceRequirement:
    ramMin: 64000
    coresMin: 4
  arv:WorkflowRunnerResources:
    ramMin: 64000
    coresMin: 4
inputs:
  directory:
    type: Directory
    inputBinding:
      position: 1
      prefix: -dir
  version:
    type: int
    inputBinding:
      position: 2
      prefix: -version
  temporaryText:
    type: File
    inputBinding:
      position: 3
      prefix: -text
  genomes:
    type: Directory[]
    inputBinding:
      position: 4

outputs:
  directoryToMerge:
    type: Directory
    outputBinding:
      glob: "."
