$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: Workflow
requirements:
  arv:RunInSingleContainer: {}
hints:
  DockerRequirement:
    dockerPull: cgivar2vcfbed
  ResourceRequirement:
    ramMin: 13000
inputs:
  cgivar: File
  sample: string
  reference: File
  cgascript:
    type: File
    default:
      class: File
      location: cgatools-mkvcf.sh
  fixscript:
    type: File
    default:
      class: File
      location: fix_vcf.py

outputs:
  vcfgz:
    type: File
    outputSource: bedtools-intersect/vcfgz
  bed:
    type: File
    outputSource: gvcf_regions/bed

steps:
  cgatools-mkvcf:
    run: cgatools-mkvcf.cwl
    in:
      cgascript: cgascript
      reference: reference
      cgivar: cgivar
      sample: sample
    out: [vcf]
  fix_vcf:
    run: fix_vcf.cwl
    in:
      fixscript: fixscript
      vcf: cgatools-mkvcf/vcf
    out: [fixedvcf]
  gvcf_regions:
    run: gvcf_regions.cwl
    in:
      vcf: fix_vcf/fixedvcf
    out: [bed]
  bedtools-intersect:
    run: bedtools-intersect.cwl
    in:
      vcf: fix_vcf/fixedvcf
      bed: gvcf_regions/bed
    out: [vcfgz]
