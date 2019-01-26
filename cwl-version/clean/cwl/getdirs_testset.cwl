$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory: Directory
  datafilenames:
    type: File
    inputBinding:
      loadContents: true

outputs:
  out1: Directory[]
  out2: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    var samplenames = [];
    var wantednames = inputs.datafilenames.contents.split('\n');
    var nlines = wantednames.length;
    for (var i = 0; i < inputs.refdirectory.listing.length; i++) {
      var name = inputs.refdirectory.listing[i];
      var type = name.class;
      var basename = name.basename;
         for (var j = 0; j < nlines; j++) {
          var wantednamej = wantednames[j];
            var result = wantednamej.indexOf(basename) > -1;
              if (result) {
              samples.push(name);
              samplenames.push(basename);
      }
     }
    }
    return {"out1": samples,"out2": samplenames};
  }

