cwlVersion: v1.1
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  chrs: string[]
  refsdir: Directory
  mapsdir: Directory
outputs:
  refs:
    type: File[]
  maps:
    type: File[]
expression: |
  ${
    var refs = [];
    var maps = [];

    for (var i = 0; i < inputs.chrs.length; i++) {
      for (var j = 0; j < inputs.refsdir.listing.length; j++) {
        var file = inputs.refsdir.listing[j];
        if (file.nameext == ".bref3" && file.basename.indexOf(inputs.chrs[i]+".") != -1) {
          refs.push(file);
        }
      }
      for (var j = 0; j < inputs.mapsdir.listing.length; j++) {
        var file = inputs.mapsdir.listing[j];
        if (file.nameext == ".map" && file.basename.indexOf(inputs.chrs[i]+".") != -1) {
          maps.push(file);
        }
      }
    }

    return {"refs": refs, "maps": maps};
  }
