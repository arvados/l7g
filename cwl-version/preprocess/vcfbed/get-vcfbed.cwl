class: ExpressionTool
cwlVersion: v1.0
inputs:
  vcfsdir: Directory
outputs:
  vcfs:
    type: File[]
    secondaryFiles: [.tbi]
  beds: File[]
  outnames: string[]
requirements:
  InlineJavascriptRequirement: {}
  cwltool:LoadListingRequirement:
    loadListing: deep_listing
expression: |
  ${
    var vcfs = [];
    var beds = [];
    var outnames = [];

    for (var i = 0; i < inputs.vcfsdir.listing.length; i++) {
      var file = inputs.vcfsdir.listing[i];
      if (file.nameext == '.gz') {
        var main = file;
        var baseName = file.nameroot.split(".")[0];
        var mainName = baseName+'.vcf.gz';
        for (var j = 0; j < inputs.vcfsdir.listing.length; j++) {
          var file = inputs.vcfsdir.listing[j];
          if (file.basename == baseName+".tbi") {
            main.secondaryFiles = [file];
          } else if (file.basename == baseName+".bed") {
            var bed = file;
          }
        }
        vcfs.push(main);
        beds.push(bed);
        outnames.push(mainName);
      }
    }
    return {"vcfs": vcfs, "beds": beds, "outnames": outnames};
  }
