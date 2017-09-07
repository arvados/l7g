class: ExpressionTool
cwlVersion: v1.0
inputs:
  pathmax: string
outputs:
  out1: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var values = [];
    var value = "0000"

    for (var i = 0; i <= inputs.pathmax; i++) {
      
      if      (i >= 0    && i <= 15)    { value = "000" + i.toString(16); }
      else if (i >= 16   && i <= 255)   { value = "00"  + i.toString(16); }
      else if (i >= 256  && i <= 4095)  { value = "0"   + i.toString(16); }
      else if (i >= 4096 && i <= 65535) { value =         i.toString(16); }          
      values.push(value)
    }
    return {"out1": values};
  } 
