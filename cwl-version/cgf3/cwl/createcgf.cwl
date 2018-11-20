$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
label: Start Arvados docker. Process and create cgf files from fastj files.
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: InlineJavascriptRequirement
  - class: ResourceRequirement
    ramMin: 10000 
    coresMin: 2 
hints:
  arv:RuntimeConstraints:
    keep_cache: 1046
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

