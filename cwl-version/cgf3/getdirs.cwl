$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
label: Create list of directories to process
requirements:
  InlineJavascriptRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  fjdir:
    type: Directory
    label: Input directory of FastJs
outputs:
  fjdirs:
    type: Directory[]
    label: Array of directories
expression: |
  ${
    var samples = [];
    for (var i = 0; i < inputs.fjdir.listing.length; i++) {
      var name = inputs.fjdir.listing[i];
      var type = name.class;
      if (type === 'Directory') {
        samples.push(name);
      }
    }
    return {"fjdirs": samples};
  }
