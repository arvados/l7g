class: ExpressionTool
cwlVersion: v1.0
inputs:
  pathmin: string
  pathmax: string
outputs:
  out1: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var values = [];
 
    var pmin = parseInt(inputs.pathmin)
    var pmax = parseInt(inputs.pathmax)

    for (var i = pmin; i <= pmax; i++) {

      var ival = i.toString()
      values.push(ival)

    }
    return {"out1": values};
  } 
