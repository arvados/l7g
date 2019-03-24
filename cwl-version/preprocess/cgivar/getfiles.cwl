class: ExpressionTool
cwlVersion: v1.0
inputs:
  dir: Directory
outputs:
  cgivars: File[]
  samples: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var cgivars = [];
    var samples = [];
    for (var i = 0; i < inputs.dir.listing.length; i++) {
      var file = inputs.dir.listing[i];
      if (file.nameext == ".bz2") {
        cgivars.push(file);
        var sample = file.basename.split(".")[0];
        samples.push(sample);
      }
    }
    return {"cgivars": cgivars, "samples": samples};
  }
