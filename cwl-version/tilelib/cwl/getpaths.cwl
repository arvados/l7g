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
    var value = "0000"

    for (var i = inputs.pathmin; i <= inputs.pathmax; i++) {
      var ival = parseInt(i);
      
      if      (ival >= 0    && ival <= 15)    { value = "000" + ival.toString(16); }
      else if (ival >= 16   && ival <= 255)   { value = "00"  + ival.toString(16); }
      else if (ival >= 256  && ival <= 4095)  { value = "0"   + ival.toString(16); }
      else if (ival >= 4096 && ival <= 65535) { value =         ival.toString(16); }          
      values.push(value)
    }
    return {"out1": values};
  } 
