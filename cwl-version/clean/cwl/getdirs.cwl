$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory:
    type: Directory
    label: Location of gVCFs to clean
outputs:
  out1
    type: Directory[]
    label: Array of directories
  out2
    type: string[]
    label: List of filename directories

requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    var samplenames = [];
    for (var i = 0; i < inputs.refdirectory.listing.length; i++) {
      var name = inputs.refdirectory.listing[i];
      var type = name.class;
      var basename = name.basename;
       if (type === 'Directory') {
              samples.push(name);
              samplenames.push(basename);
            }
    }
    return {"out1": samples,"out2": samplenames};
  }
