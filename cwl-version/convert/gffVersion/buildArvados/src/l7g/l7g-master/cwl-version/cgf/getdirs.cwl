class: ExpressionTool
cwlVersion: v1.0
inputs:
  refdirectory: Directory
outputs:
  out1: Directory[] 
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    for (var i = 0; i < inputs.refdirectory.listing.length; i++) {
      var name = inputs.refdirectory.listing[i];
      var type = name.class;
       if (type === 'Directory') {
            samples.push(name)
          }
    }
    return {"out1": samples};
  } 
