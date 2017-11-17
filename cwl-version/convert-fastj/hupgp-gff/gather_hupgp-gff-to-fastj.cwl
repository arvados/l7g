class: ExpressionTool
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}

inputs:
  indirs: Directory[]
outputs:
  out: Directory
expression: |
  ${
    //var gathered_dirs = new Directory;
    var gathered_dirs = { "class" : "Directory", "basename": "output", "listing" : [] };
    for (var i=0; i<inputs.indirs.length; i++) {
      var ele = inputs.indirs[i];

      //gathered_dirs.listing.concat(ele.listing);

      gathered_dirs.listing.push(ele.listing[0]);

      // gives output with keep:* subdirectories.
      //
      //gathered_dirs.listing.push(inputs.indirs[i]);

      // gives a 'no get method' error
      //
      //gathered_dirs.listing.push(inputs.indirs[i].listing);

    }
    var x = JSON.stringify(gathered_dirs);
    //process.stdout.write(x);
    return { "out": gathered_dirs };
  }

