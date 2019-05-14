$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Filters gVCFs by a specified quality cutoff and cleans
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ScatterFeatureRequirement: {}
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096

inputs:
  gvcfdir:
    type: Directory
    label: Input gVCF directory
  cutoff:
    type: int
    label: Filtering cutoff threshold

outputs:
  filteredcleangvcfs:
    type: File[]
    label: Filtered clean gVCFs
    outputSource: filtercleangvcf/filteredcleangvcf

steps:
  getfiles:
    run: getfiles.cwl
    in:
      gvcfdir: gvcfdir
    out: [gvcfs]

  filtercleangvcf:
    run: filtercleangvcf.cwl
    scatter: gvcf
    in: 
      gvcf: getfiles/gvcfs
      cutoff: cutoff
    out: [filteredcleangvcf]
