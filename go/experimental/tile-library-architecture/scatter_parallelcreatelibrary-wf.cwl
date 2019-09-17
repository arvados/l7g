#!/usr/bin/env cwl-runner
cwlVersion: v1.0
class: Workflow

requirements:
  ScatterFeatureRequirement: {}

hints:
  ResourceRequirement:
    ramMin: 64000
    cpuMin: 4

inputs:
  directory:
    type: string
    inputBinding:
      position: 1
      prefix: -dir
  version:
    type: int
    inputBinding:
      position: 2
      prefix: -version
  genomesList:
    type:
      type: array
      items:
        type: array
        items: string
    inputBinding:
      position: 3

steps:
  mergefastj:
    run: mergefastj.cwl
    scatter: genomes
    in:
      version: version
      genomes: genomesList
    out: [directoryToMerge]
  mergeintermediatelibraries:
    run: mergeintermediatelibraries.cwl
    in:
      version: version
      directory: directory
      libraries: createnewlibrary/[directoryToMerge]
    out: []
outputs: []