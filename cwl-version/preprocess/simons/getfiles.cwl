cwlVersion: v1.0
class: ExpressionTool
label: Create list of VCFs and sample names
inputs:
  dir:
    type: Directory
    label: Input directory of VCFs
outputs:
  vcfs:
    type: File[]
    label: Output VCFs
  samples:
    type: string[]
    label: Sample names of VCFs
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var vcfs = [];
    var samples = [];
    for (var i = 0; i < inputs.dir.listing.length; i++) {
      var file = inputs.dir.listing[i];
      if (file.nameext == ".gz") {
        vcfs.push(file);
        var sample = file.basename.split(".")[0];
        samples.push(sample);
      }
    }
    return {"vcfs": vcfs, "samples": samples};
  }
