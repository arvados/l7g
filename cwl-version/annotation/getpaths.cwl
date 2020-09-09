cwlVersion: v1.1
class: ExpressionTool
label: Create list of paths from tile library
requirements:
  InlineJavascriptRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  tilelib:
    type: Directory
    label: Input tile library
outputs:
  pathstrs:
    type: string[]
    label: Output path strings
expression: |
  ${
    var pathstrs = [];
    for (var i = 0; i < inputs.tilelib.listing.length; i++) {
      var file = inputs.tilelib.listing[i];
      if (file.nameext == ".gz") {
        var pathstr = file.basename.split(".")[0];
        pathstrs.push(pathstr);
      }
    }
    pathstrs.sort();
    return {"pathstrs": pathstrs};
  }
