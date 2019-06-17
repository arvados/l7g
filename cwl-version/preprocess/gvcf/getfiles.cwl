$namespaces:
  cwltool: "http://commonwl.org/cwltool#"
class: ExpressionTool
label: Create list of gVCFs from directory
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  gvcfdir:
    type: Directory
    label: Directory of input gVCFs
outputs:
  gvcfs:
    type: File[]
    label: Array of gvcfs
    secondaryFiles: [.tbi]
expression: |
  ${
    var gvcfs = [];
    for (var i = 0; i < inputs.gvcfdir.listing.length; i++) {
      var file = inputs.gvcfdir.listing[i];
      if (file.nameext == '.gz') {
        var main = file;
        for (var j = 0; j < inputs.gvcfdir.listing.length; j++) {
          var file = inputs.gvcfdir.listing[j];
          if (file.basename == main.basename+".tbi") {
            main.secondaryFiles = [file];
          }
        }
        gvcfs.push(main);
      }
    }
    return {"gvcfs": gvcfs};
  }
