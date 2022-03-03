cwlVersion: v1.1
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  sample: string
  chrs: string[]
  vcfdir: Directory
  gnomaddir: Directory
outputs:
  samples: string[]
  vcfs: File[]
  gnomads:
    type: File[]
    secondaryFiles: [.csi]
expression: |
  ${
    var samples = [];
    var vcfs = [];
    var gnomads = [];

    for (var i = 0; i < inputs.chrs.length; i++) {
      var chr = inputs.chrs[i];
      var sample = inputs.sample+"."+chr;
      for (var j = 0; j < inputs.vcfdir.listing.length; j++) {
        var file = inputs.vcfdir.listing[j];
        if (file.basename.includes("."+chr+".")) {
          var vcf = file;
          break;
        }
      }
      for (var j = 0; j < inputs.gnomaddir.listing.length; j++) {
        var file = inputs.gnomaddir.listing[j];
        if (file.basename.includes("."+chr+".")) {
          var gnomad = file;
          break;
        }
      }
      for (var j = 0; j < inputs.gnomaddir.listing.length; j++) {
        var file = inputs.gnomaddir.listing[j];
        if (file.basename == gnomad.basename+".csi") {
          gnomad.secondaryFiles = [file];
          break;
        }
      }
      samples.push(sample);
      vcfs.push(vcf);
      gnomads.push(gnomad);
    }

    return {"samples": samples, "vcfs": vcfs, "gnomads": gnomads};
  }
