cwlVersion: v1.0
class: Workflow
label: Scatter to process portable VCFs
requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: vcfutil

inputs:
  vcfsdir:
    type: Directory
    label: Input directory of VCFs
  header:
    type: File
    label: Header file
    default:
      class: File
      location: header
  sdf:
    type: Directory
    label: RTG reference directory
  cleanvcf:
    type: File
    label: Code that cleans VCFs
    default:
      class: File
      location: ../gvcf/src/cleanvcf.py

outputs:
  processedvcfgzs:
    type: File[]
    label: Processed VCFs
    outputSource: preprocess-portablevcf-wf/processedvcfgz

steps:
  getfiles:
    run: getfiles.cwl
    in:
      dir: vcfsdir
    out: [vcfgzs]
  preprocess-portablevcf-wf:
    run: preprocess-portablevcf-wf.cwl
    scatter: vcfgz
    in:
      vcfgz: getfiles/vcfgzs
      header: header
      sdf: sdf
      cleanvcf: cleanvcf
    out: [processedvcfgz, summary]
  cat:
    run: cat.cwl
    in:
      txts: preprocess-portablevcf-wf/summary
    out: [cattxt]
