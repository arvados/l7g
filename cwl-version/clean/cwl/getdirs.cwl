$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
cwlVersion: v1.0
label: Create list of gVCFs directories to clean
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdirectory:
    type: Directory
    label: Input directory of gVCFs
outputs:
  out1:
    type: Directory[]
    label: Array of directories containing gVCFs
  out2:
    type: string[]
    label: Array of directory names

requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var samples = [];
    var samplenames = [];
    for (var i = 0; i < inputs.refdirectory.listing.length; i++) {
      var name = inputs.refdirectory.listing[i];
      var type = name.class;
      var basename = name.basename;
       if (type === 'Directory') {
              samples.push(name);
              samplenames.push(basename);
            }
    }
    return {"out1": samples,"out2": samplenames};
  }
