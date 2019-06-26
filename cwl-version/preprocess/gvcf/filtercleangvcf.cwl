$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Filters gVCFs by a specified quality cutoff and cleans
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  bashscript:
    type: File
    label: Master script to control filtering
    default:
      class: File
      location: src/filtercleanCWL.sh
  gvcf:
    type: File
    label: Input gVCF file
    secondaryFiles: [.tbi]
  filtergvcf:
    type: File
    label: Code that filters gVCFs
    default:
      class: File
      location: src/filter-gvcf
  cutoff:
    type: int
    label: Filtering cutoff threshold
  cleanvcf:
    type: File
    label: Code that cleans gVCFs
    default:
      class: File
      location: src/cleanvcf.py
outputs:
  filteredcleangvcf:
    type: File
    label: Filtered and clean gVCF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: bash
arguments:
  - $(inputs.bashscript)
  - $(inputs.gvcf)
  - $(inputs.filtergvcf)
  - $(inputs.cutoff)
  - $(inputs.cleanvcf)
