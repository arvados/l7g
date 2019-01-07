class: ExpressionTool
cwlVersion: v1.0
inputs:
  pathmin: string
  pathmax: string
  nchunks: string
outputs:
  out1: string[]
  out2: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${

    var imax = parseInt(inputs.pathmax);
    var imin = parseInt(inputs.pathmin);
    var chunk_size = parseInt(inputs.nchunks);

    var index = 0;
    var myArray = [];
    var tempArray = [];
    var maxArray = [];
    var minArray = [];

    for (var ival = imin; ival <= imax; ival++) {
      var value = ival;
      myArray.push(value)
    }

    var arrayLength = myArray.length;

    for (index = 0; index < arrayLength; index += chunk_size) {
       var myChunk = myArray.slice(index, index+chunk_size);
       var minval = myChunk[0];
       var minvalstr = minval.toString();
       var maxval = myChunk[myChunk.length-1];
       var maxvalstr = maxval.toString();
       maxArray.push(maxvalstr);
       minArray.push(minvalstr);
    }

    return {"out1": minArray, "out2": maxArray};
  }
