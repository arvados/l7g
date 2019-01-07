$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow 
requirements:
  - class: DockerRequirement
    dockerPull: arvados/l7g
  - class: ScatterFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 16384 
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  pathmin: string
  pathmax: string
  nchunks: string
  datadir: Directory
  tagset: File

outputs:
  out1:
    type:
      type: array
      items:
        type: array
        items: File
    outputSource: step2/out1

steps:
  step1:
    run: getpaths_chunk.cwl
    in:
      pathmin: pathmin
      pathmax: pathmax
      nchunks: nchunks
    out: [out1,out2]

  step2:
     scatter: [tilepathmin,tilepathmax]
     scatterMethod: dotproduct
     in:
       tilepathmin: step1/out1
       tilepathmax: step1/out2
       datadir: datadir
       tagset: tagset
     run: createsglf_chunkv2.cwl
     out: [out1]
