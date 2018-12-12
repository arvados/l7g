$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: ExpressionTool
label: Create list of gVCF to be handled
requirements:
  - class: InlineJavascriptRequirement

hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  datafilenames:
    label: Directories of gVCF chromosome files
    type:
      type: array
      items: File
      inputBinding:
        loadContents: true
  datafilepdh:
    label: List of portable data hashes (pdh) in Arvados
    type:
      type: array
      items: File
      inputBinding:
        loadContents: true

outputs:
  fileprefix:
    type: string[]
    label: Array of directory names
  collectiondir:
    type: Directory[]
    label: Array of directories containing gVCFs

expression: |
   ${
     var fileprefix=[];
     var collectiondir=[];
     var ssdir=[];
     for (var jj = 0; jj < inputs.datafilenames.length; jj++) {
      var names=inputs.datafilenames[jj].contents.split('\n');
      var nblines1=names.length;
      for (var j = 0; j < nblines1-1; j++) {
       var nn=names[j];
       fileprefix.push(nn);
       }
      }

     for (var ii = 0; ii < inputs.datafilepdh.length; ii++) {
       var pdhs=inputs.datafilepdh[ii].contents.split('\n');
       var nblines2=pdhs.length;
       for (var i = 0; i < nblines2-1; i++) {
         var ss=pdhs[i];
         var ssdir={"class": "Directory", "location": "keep:" + ss};
         collectiondir.push(ssdir);
         }
       }

     return {"fileprefix": fileprefix,"collectiondir":collectiondir};
     }
