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
    ramMin: 57344 
    ramMax: 114688 
    coresMin: 8
    coresMax: 8
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
  fjdir:
    type: Directory 
    inputBinding:
      position: 2
  cgft: 
    type: File 
    inputBinding:
      position: 3 
  fjt:
    type: File
    inputBinding:
      position: 4
  cglf:
    type: Directory
    inputBinding:
      position: 5 
outputs:    
  out1: 
    type: File 
    outputBinding:
      glob: "data/*.cgf"

