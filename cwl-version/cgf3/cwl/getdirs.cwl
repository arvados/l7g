$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory: Directory
outputs:
  out1: Directory[] 
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
