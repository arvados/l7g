cwlVersion: v1.1
class: ExpressionTool
label: Create list of VCFs and sample names
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  vcfsdir: Directory
outputs:
  vcfs:
    type: File[]
    secondaryFiles: [.tbi]
  samples: string[]
requirements:
  InlineJavascriptRequirement: {}
expression: |
  ${
    var vcfs = [];
    var samples = [];

    for (var i = 0; i < inputs.vcfsdir.listing.length; i++) {
      var file = inputs.vcfsdir.listing[i];
      if (file.nameext == '.gz') {
        var sample = file.nameroot.split('.').slice(0,-1).join('.');
        var main = file;
        for (var j = 0; j < inputs.vcfsdir.listing.length; j++) {
          var file = inputs.vcfsdir.listing[j];
          if (file.basename == main.basename+".tbi") {
            main.secondaryFiles = [file];
          }
        }
        vcfs.push(main);
        samples.push(sample);
      }
    }
    return {"vcfs": vcfs, "samples": samples};
  }
