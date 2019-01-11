class: ExpressionTool
cwlVersion: v1.0
label: Create list of gff directories to process
inputs:
  refdirectory:
    type: Directory
    label: Location of gff to convert
outputs:
  out1:
    type: File[]
    label: Array of gffs
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    for (var i = 0; i < inputs.refdirectory.listing.length; i++) {
      var file = inputs.refdirectory.listing[i];
      var filename = file.basename;
      var field = filename.split('.').pop();
          if (field === 'gz') {
            samples.push(file)
          }
    }
    return {"out1": samples};
  }
