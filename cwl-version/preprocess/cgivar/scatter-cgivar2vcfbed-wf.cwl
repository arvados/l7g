$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: Workflow
requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: process-cgi
inputs:
  cgivarsdir: Directory
  reference: File

outputs:
  vcfgzs:
    type: File[]
    outputSource: cgivar2vcfbed-wf/vcfgz
  beds:
    type: File[]
    outputSource: cgivar2vcfbed-wf/bed

steps:
  getfiles:
    run: getfiles.cwl
    in:
      dir: cgivarsdir
    out: [cgivars, samples]
  cgivar2vcfbed-wf:
    requirements:
      arv:RunInSingleContainer: {}
    run: cgivar2vcfbed-wf.cwl
    scatter: [cgivar, sample]
    scatterMethod: dotproduct
    in:
      cgivar: getfiles/cgivars
      sample: getfiles/samples
      reference: reference
    out: [vcfgz, bed]
