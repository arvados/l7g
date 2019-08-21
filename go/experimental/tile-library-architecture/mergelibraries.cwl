#!/usr/bin/env cwl-runner
cwlVersion: v1.0
class: CommandLineTool
baseCommand: ./mergelibraries

hints:
  ResourceRequirement:
    coresMin: 16
    ramMin: 64000
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
  libraries:
    type: Directory[]
    inputBinding:
      position: 3

outputs: []
