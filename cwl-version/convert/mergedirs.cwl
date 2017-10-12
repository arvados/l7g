$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: ExpressionTool
requirements:
  - class: InlineJavascriptRequirement
hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  refdir1: Directory
  refdir2: Directory
outputs:
  out1: Directory[]
expression: |
  ${
    var alldirs = [];
    for (var i = 0; i < inputs.refdir1.listing.length; i++) {
      var content = inputs.refdir1.listing[i];
         if (content.class === 'Directory') {
            alldirs.push(content)
         }
    }

   for (var j = 0; j < inputs.refdir2.listing.length; j++) {
      var content = inputs.refdir2.listing[j];
         if (content.class === 'Directory') {
            alldirs.push(content)
         }
    }

    return {"out1": alldirs};
   } 
