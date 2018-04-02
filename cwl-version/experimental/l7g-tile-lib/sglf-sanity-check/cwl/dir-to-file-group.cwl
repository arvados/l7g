class: ExpressionTool
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}

inputs:
  inputdir: Directory

outputs:
  out: File[]

expression: |
  ${
    var dir_array = [];

    var n = inputs.inputdir.listing.length;
    var m = Math.ceil( n / 20 );
    for (var i=0; i<m; i++) {
      dir_array.push( {"out" : [] } );
    }

    var samples = [];
    for (var i = 0; i < inputs.inputdir.listing.length; i++) {
      var file = inputs.inputdir.listing[i];
      samples.push(file);

      var k = Math.floor( i / m );

      dir_array[k].push(file);

    }
    return {"out": samples};
  }
