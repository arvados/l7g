class: ExpressionTool
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}
inputs:
  gffdir: Directory
outputs:
  gffs: File[]
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
