$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Preprocess VCF and BED files to create a collection of gVCF files
requirements:
  DockerRequirement:
    dockerPull: l7g/preprocess-vcfbed
  ScatterFeatureRequirement: {}

inputs:
  vcfsdir:
    type: Directory
    label: Directory of VCF, BED and index files
  ref:
    type: File
    label: Reference FASTA file
  bedfile:
    type: File?
    label: Optional BED to scatter over if not included in vcfsdir

outputs:
  result:
    type: File[]
    label: gVCFs and index files
    outputSource: vcfbed2gvcf/result

steps:
  get-vcfbed:
    run: get-vcfbed.cwl
    in:
      vcfsdir: vcfsdir
      bedfile: bedfile
    out: [vcfs, beds, outnames]
  vcfbed2gvcf:
    run: vcfbed2gvcf.cwl
    scatter: [vcf, bed, outname]
    scatterMethod: dotproduct
    in:
      vcf: get-vcfbed/vcfs
      bed: get-vcfbed/beds
      ref: ref
      outname: get-vcfbed/outnames
    out: [result]
