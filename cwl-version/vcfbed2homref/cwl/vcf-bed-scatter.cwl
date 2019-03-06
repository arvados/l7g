class: ExpressionTool
cwlVersion: v1.0
inputs:
  vcfsdir: Directory
outputs:
  gvcfFns:
    type: File[]
    secondaryFiles: [.tbi]
  bedFns: File[]
  outNames: string[]
requirements:
  InlineJavascriptRequirement: {}
  cwltool:LoadListingRequirement:
    loadListing: deep_listing
expression: |
  ${
    var gvcfFns = [];
    var bedFns = [];
    var outNames = [];

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
        gvcfFns.push(main);
        bedFns.push(bed);
        outNames.push(mainName);

      }
    }

    return {"gvcfFns": gvcfFns, "bedFns": bedFns, "outNames": outNames};
  }
