$namespaces:
  arv: "http://arvados.org/cwl#"
cwlVersion: v1.0
class: Workflow
label: Convert CGIVAR to VCF and BED region
requirements:
  arv:RunInSingleContainer: {}
hints:
  DockerRequirement:
    dockerPull: cgivar2vcfbed
  ResourceRequirement:
    ramMin: 13000
inputs:
  cgivar:
    type: File
    label: Input CGIVAR
  sample:
    type: string
    label: Sample name
  reference:
    type: File
    label: CRR reference used for cgatools
  cgascript:
    type: File
    label: Script invoking cgatools
    default:
      class: File
      location: cgatools-mkvcf.sh
  fixscript:
    type: File
    label: Script to fix VCF
    default:
      class: File
      location: fix_vcf.py

outputs:
  vcfgz:
    type: File
    label: Output VCF
    outputSource: bedtools-intersect/vcfgz
  bed:
    type: File
    label: BED region VCF
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
