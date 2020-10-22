$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.1
class: CommandLineTool
label: Convert CGF to VCF
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
  ResourceRequirement:
    coresMin: 2
    ramMin: 8000
hints:
  arv:RuntimeConstraints:
    keep_cache: 4000
inputs:
  bashscript:
    type: File
    label: Shel script to convert CGF to VCF
    default:
      class: File
      location: ../../tools/cgf2vcf/cgf2vcf.sh
  makeheadergenomebed:
    type: File
    label: Python script to make VCF header and genome BED
    default:
      class: File
      location: ../../tools/cgf2vcf/makeheadergenomebed.py
  cgf2vcf:
    type: File
    label: Python script to convert CGF to VCF
    default:
      class: File
      location: ../../tools/cgf2vcf/cgf2vcf.py
  steps2bed:
    type: File
    label: Python script to convert steps to BED
    default:
      class: File
      location: ../../tools/cgf2vcf/steps2bed.py
  assembly:
    type: File
    label: Compressed assembly fixed width file
    secondaryFiles: [^.fwi, .gzi]
  annotationlib:
    type: Directory
    label: Input annotation library
  cgf:
    type: File
    label: Input CGF file
outputs:
  vcf:
    type: File
    label: VCF converted from CGF
    outputBinding:
      glob: "*vcf.gz"
    secondaryFiles: [.tbi]
  bed:
    type: File
    label: Coverage BED
    outputBinding:
      glob: "*bed"
arguments:
  - $(inputs.bashscript)
  - $(inputs.makeheadergenomebed)
  - $(inputs.cgf2vcf)
  - $(inputs.steps2bed)
  - $(inputs.assembly)
  - $(inputs.annotationlib)
  - $(inputs.cgf)
