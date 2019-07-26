$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Preprocess VCF and BED files to create a collection of gVCF files
requirements:
  - class: DockerRequirement
    dockerPull: l7g/preprocess-vcfbed
  - class: ResourceRequirement
    coresMin: 1
    ramMin: 12000
  - class: ScatterFeatureRequirement

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
    label: Directory containing gVCF and index files
    outputSource: vcfbed2gvcf/result

steps:
  get-vcfbed:
    run: get-vcfbed.cwl
    in:
      vcfsdir: vcfsdir
      bedfile: bedfile
    out: [vcfs, beds, outnames]
  sort-vcf:
    run: sort-vcf.cwl
    scatter: vcf
    in:
      vcf: get-vcfbed/vcfs
    out: [sortedvcf]
  sort-bed:
    scatter: bed
    run: sort-bed.cwl
    in:
      bed: get-vcfbed/beds
    out: [sortedbed]
  intersect-vcfbed:
    run: intersect-vcfbed.cwl
    scatter: [vcf, bed]
    scatterMethod: dotproduct
    in:
      vcf: sort-vcf/sortedvcf
      bed: sort-bed/sortedbed
    out: [intersectedvcf]
  vcfbed2gvcf:
    run: vcfbed2gvcf.cwl
    scatter: [vcf, bed, outname]
    scatterMethod: dotproduct
    in:
      vcf: intersect-vcfbed/intersectedvcf
      bed: sort-bed/sortedbed
      ref: ref
      outname: get-vcfbed/outnames
    out: [result]
