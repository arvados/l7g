cwlVersion: v1.1
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
inputs:
  dirs:
    type:
      type: array
      items: Directory
  batchsize:
    type: int
outputs:
  batches:
    type:
      type: array
      items:
        type: array
        items: Directory
expression: |
  ${
    var batches = [];
    for (var i = 0; i < inputs.dirs.length; i+=inputs.batchsize) {
      var batch = inputs.dirs.slice(i, i+inputs.batchsize);
      batches.push(batch);
    }
    return {"batches": batches};
  }
