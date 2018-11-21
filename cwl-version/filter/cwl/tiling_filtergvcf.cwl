$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Filters GVCFs by some quality cutoff
doc: |
    This workflow takes in RAW gVCFs, and using the cutoff integer as a quality
    cutoff filters out gVCFs in the given path that do not meat that cutoff.
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  datafilenames: File[]
  datafilepdh: File[]
  bashscript: File
  filter_gvcf: File
  cutoff: string
outputs:
  out1:
    type: Directory[]
    outputSource: step2/out1
    label: Output directory of filterd gVCFs
steps:
  step1:
    run: getCollections.cwl
    in: 
      datafilenames: datafilenames
      datafilepdh: datafilepdh
    out: [fileprefix,collectiondir]

  step2:
    scatter: [gffPrefix,gffDir] 
    scatterMethod: dotproduct
    in: 
      bashscript: bashscript
      gffDir: step1/collectiondir
      gffPrefix: step1/fileprefix 
      filter_gvcf: filter_gvcf
      cutoff: cutoff
    run: filter.cwl
    out: [out1]
