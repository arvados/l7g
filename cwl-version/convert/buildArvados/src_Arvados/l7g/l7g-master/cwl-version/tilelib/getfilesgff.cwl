class: ExpressionTool
cwlVersion: v1.0
inputs:
  refdirectory: Directory
outputs:
  out1: File[]
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
