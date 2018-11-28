$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create a tile library (SGLF) for a given set of FASTJ files
doc: |
    Creates the tile library set (SGLF) in Compact Genome Format (cgf)
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
  pathmin
    type: string
  pathmax
    type: string
  nchunks
    type: string
  bashscript
    type: File
  fjcsv2sglf
    type: File
  datadir
    type: Directory
  fjt
    type: File
  tagset
    type: File

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
