cwlVersion: v1.2
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  fas:
    type:
      type: array
      items:
        type: array
        items: File
outputs:
  fastadirs: Directory[]
expression: |
  ${
    var fastadirs = [];
    for (var i = 0; i < inputs.fas.length; i+=100) {
      var fastadir = {"class": "Directory",
                      "basename": "dir"+String(i/100),
                      "listing": []};
      for (var j = i; j < Math.min(i+100, inputs.fas.length); j+=1) {
        fastadir.listing.push(inputs.fas[j][0]);
        fastadir.listing.push(inputs.fas[j][1]);
      }
      fastadirs.push(fastadir);
    }
    return {"fastadirs": fastadirs};
  }
