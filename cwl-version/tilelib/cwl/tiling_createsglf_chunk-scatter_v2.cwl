$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create a tile library (SGLF) for a given set of FastJ files
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
  - class: ScatterFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 16384
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  pathmin:
    type: string
    label: Path to start at in the tile library
  pathmax:
    type: string
    label: Last/Maximum tile in library
  nchunks:
    type: string
    label: Number of chunks to scatter
  bashscript:
    type: File
    label: Script that iterates over the FastJ to create paths
  fjcsv2sglf:
    type: File
    label: Tool to create tile library
  datadir:
    type: Directory
    label: Directory of FastJ files
  fjt:
    type: File
    label: Tool to manipulate FastJ (text) files
  tagset:
    type: File
    label: Compressed tagset in FASTA format
outputs:
  out1:
    type:
      type: array
      items:
        type: array
        items: File
    label: Output tile library
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
     scatter: [tilepathmin, tilepathmax]
     scatterMethod: dotproduct
     in:
       bashscript: bashscript
       tilepathmin: step1/out1
       tilepathmax: step1/out2
       fjcsv2sglf: fjcsv2sglf
       datadir: datadir
       fjt: fjt
       tagset: tagset
     run: createsglf_chunkv2.cwl
     out: [out1]
