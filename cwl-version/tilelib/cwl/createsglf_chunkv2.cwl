$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    ramMin: 100000
    coresMin: 16
    coresMax: 16
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

baseCommand: bash
inputs:
  bashscript:
    type: File
    inputBinding:
      position: 1
  tilepathmin:
    type: string
    inputBinding:
      position: 2
  tilepathmax:
    type: string
    inputBinding:
      position: 3
  fjcsv2sglf: 
    type: File 
    inputBinding:
      position: 4
  datadir: 
    type: Directory 
    inputBinding:
      position: 5 
  fjt:
    type: File
    inputBinding:
      position: 6 
  tagset:
    type: File
    inputBinding:
      position: 7 

outputs:
  out1:
    type: File[]
    outputBinding:
      glob: "lib/*sglf.gz*"

