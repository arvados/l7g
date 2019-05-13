cwlVersion: v1.0
class: ExpressionTool
label: Create list of CGIVARs to process
inputs:
  dir:
    type: Directory
    label: Input directory of CGIVARs
outputs:
  cgivars:
    type: File[]
    label: Output CGIVARs
  samples:
    type: string[]
    label: Sample names of CGIVARs
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
