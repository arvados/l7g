$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
label: Create list of directories to process
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory: Directory
outputs:
  out1: Directory[]
  out2: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    var samplenames = [];
    var n = 0
    var i = 0
   
    do {
      var name = inputs.refdirectory.listing[i];
      var type = name.class;
      var basename = name.basename; 
      i +=1; 
        if (type === 'Directory') {
              n += 1;
              samples.push(name);
              samplenames.push(basename);
            }
      } while (n < 20);
    return {"out1": samples,"out2": samplenames};
  } 
