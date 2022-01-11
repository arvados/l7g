cwlVersion: v1.1
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
inputs:
  nestedarray:
    type:
      type: array
      items:
        type: array
        items: Directory
outputs:
  flattenedarray:
    type:
      type: array
      items: Directory
expression: |
  ${
    var flattenedarray = [];
    for (var i = 0; i < inputs.nestedarray.length; i++) {
      for (var j = 0; j < inputs.nestedarray[i].length; j++) {
        flattenedarray.push(inputs.nestedarray[i][j]);
      }
    }
    return {"flattenedarray": flattenedarray};
  }
