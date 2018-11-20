class: ExpressionTool
cwlVersion: v1.0
label: Collect list of validated cgfs, return JSON bundle
requirements:
  InlineJavascriptRequirement: {}

inputs:
  infiles: File[]
outputs:
  out: Directory
expression: |
  ${
    var gathered_dirs = { "class" : "Directory", "basename": "output", "listing" : [] };
    for (var i=0; i<inputs.infiles.length; i++) {

      var ele = inputs.infiles[i];
      //var d = { "class" : "Directory", "basename" : String(i) + "-" + String(ele.basename), "listing": [] };

      // The first element in the listing is the file we care about?
      //
      //gathered_dirs.listing.push(ele.listing[0]);
      gathered_dirs.listing.push(ele);

    }
    //var x = JSON.stringify(gathered_dirs);
    return { "out": gathered_dirs };
  }

