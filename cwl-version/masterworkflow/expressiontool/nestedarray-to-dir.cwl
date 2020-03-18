$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.1
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  nestedarr:
    type:
      type: array
      items:
        type: array
        items: [File, Directory]
  dirname:
    type: string
outputs:
  dir: Directory
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var dir = {"class": "Directory",
               "basename": inputs.dirname,
               "listing": []};
    for (var i = 0; i < inputs.nestedarr.length; i++) {
      for (var j = 0; j < inputs.nestedarr[i].length; j++) {
        dir.listing.push(inputs.nestedarr[i][j]);
      }
    }
    return {"dir": dir};
  }
