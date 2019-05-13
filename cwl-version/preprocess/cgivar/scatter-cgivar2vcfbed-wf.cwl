$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: Workflow
label: Scatter to convert CGIVARs to VCFs and BEDs
requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}
hints:
  DockerRequirement:
    dockerPull: process-cgi
inputs:
  cgivarsdir:
    type: Directory
    label: Input directory of CGIVARs
  reference:
    type: File
    label: CRR reference used for cgatools

outputs:
  vcfgzs:
    type: File[]
    label: Output VCFs
    outputSource: cgivar2vcfbed-wf/vcfgz
  beds:
    type: File[]
    label: Output BEDs
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
