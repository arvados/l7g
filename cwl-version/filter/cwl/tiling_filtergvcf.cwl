$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Filters gVCFs by a specified quality cutoff
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  datafilenames:
    type: File[]
    label: Files listing input gVCFs
  datafilepdh:
    type: File[]
    label: Files listing input portable data hashes
  cutoff:
    type: string
    label: Filtering cutoff threshold

outputs:
  out1:
    type: Directory[]
    label: Directory of filtered gVCFs
    outputSource: step2/out1

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
      gffDir: step1/collectiondir
      gffPrefix: step1/fileprefix
      cutoff: cutoff
    run: filter.cwl
    out: [out1]
