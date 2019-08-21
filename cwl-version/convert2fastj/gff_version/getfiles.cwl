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
    secondaryFiles: [.tbi]
expression: |
  ${
    var gffs = [];
    for (var i = 0; i < inputs.gffdir.listing.length; i++) {
      var file = inputs.gffdir.listing[i];
      if (file.nameext == '.gz') {
        var main = file;
        for (var j = 0; j < inputs.gffdir.listing.length; j++) {
          var file = inputs.gffdir.listing[j];
          if (file.basename == main.basename+".tbi") {
            main.secondaryFiles = [file];
          }
        }
        gffs.push(main);
      }
    }
    return {"gffs": gffs};
  }
