class: ExpressionTool
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}

inputs:
  idirs: Directory[]
outputs:
  result: Directory
expression: |
  ${
    var gathered_dirs = { "class" : "Directory", "basename": ".", "listing" : [] };
    for (var i=0; i<inputs.idirs.length; i++) {

      var ele = inputs.idirs[i];

      // The first element in the listing is the directory we care about.
      //
      gathered_dirs.listing.push(ele.listing[0]);

    }
    return { "result": gathered_dirs };
  }
