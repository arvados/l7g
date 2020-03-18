$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.1
class: Workflow
label: Create a tile library (SGLF) for a given set of FastJ files
requirements:
  ScatterFeatureRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: no_listing

inputs:
  pathmin:
    type: string
    label: Starting path in the tile library
  pathmax:
    type: string
    label: Last/Maximum path in the tile library
  nchunks:
    type: string
    label: Number of chunks to scatter
  fjdir:
    type: Directory
    label: Directory of FastJ files
  tagset:
    type: File
    label: Compressed tagset in FASTA format

outputs:
  sglfs:
    type:
      type: array
      items:
        type: array
        items: File
    label: Output tile library
    outputSource: createsglf/chunksglfs

steps:
  getpaths:
    run: getpaths.cwl
    in:
      pathmin: pathmin
      pathmax: pathmax
      nchunks: nchunks
    out: [minpaths, maxpaths]

  createsglf:
     run: createsglf.cwl
     scatter: [tilepathmin, tilepathmax]
     scatterMethod: dotproduct
     in:
       tilepathmin: getpaths/minpaths
       tilepathmax: getpaths/maxpaths
       fjdir: fjdir
       tagset: tagset
     out: [chunksglfs]
