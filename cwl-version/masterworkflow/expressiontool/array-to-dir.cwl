$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.1
hints:
  LoadListingRequirement:
    loadListing: no_listing
inputs:
  arr:
    type:
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
               "listing": inputs.arr};
    return {"dir": dir};
  }
