$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Create NumPy arrays by tile path from cgfs, merge all NumPy arrays into single array
requirements:
  DockerRequirement:
    dockerPull: arvados/l7g
hints:
  cwltool:LoadListingRequirement:
    loadListing: no_listing

inputs:
  waitsignal:
    type: Any?
    label: Wait signal to start workflow
  cgfdir:
    type: Directory
    label: Directory of compact genome format files

outputs:
  consolnpydir:
    type: Directory
    label: Output consolidated NumPy arrays
    outputSource: consolnpy/consolnpydir
  names:
    type: File
    label: File listing sample names
    outputSource: createnpy/names

steps:
  createnpy:
    run: createnpy.cwl
    in:
      cgfdir: cgfdir
    out: [npydir, names]

  consolnpy:
    run: consolnpy.cwl
    in: 
      npydir: createnpy/npydir
    out: [consolnpydir]
