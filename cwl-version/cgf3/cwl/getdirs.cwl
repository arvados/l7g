$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
label: Create list of directories to process
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory: Directory
  label: Directory of files to be analyzed
outputs:
  out1: Directory[]
  label: List of processed directories
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    for (var i = 0; i < inputs.refdirectory.listing.length; i++) {
      var name = inputs.refdirectory.listing[i];
      var type = name.class;
       if (type === 'Directory') {
            samples.push(name)
          }
    }
    return {"out1": samples};
  }
