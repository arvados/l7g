$namespaces:
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
label: Create list of GFFs from directory
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  gffdir:
    type: Directory
    label: Directory of input GFFs
outputs:
  gffs:
    type: File[]
    label: Array of GFFs
expression: |
  ${
    var gffs = [];
    for (var i = 0; i < inputs.gffdir.listing.length; i++) {
      var file = inputs.gffdir.listing[i];
      if (file.nameext == '.gz') {
        gffs.push(file);
      }
    }
    return {"gffs": gffs};
  }
