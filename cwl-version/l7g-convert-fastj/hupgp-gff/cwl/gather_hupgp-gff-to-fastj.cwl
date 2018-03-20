class: ExpressionTool
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}

inputs:
  indirs: Directory[]
outputs:
  out: Directory
expression: |
  ${
    var gathered_dirs = { "class" : "Directory", "basename": "output", "listing" : [] };
    for (var i=0; i<inputs.indirs.length; i++) {

      var ele = inputs.indirs[i];

      // The first element in the listing is the directory we care about.
      //
      gathered_dirs.listing.push(ele.listing[0]);

    }
    return { "out": gathered_dirs };
  }

