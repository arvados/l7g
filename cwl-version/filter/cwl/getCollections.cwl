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
  datafilenames:
   type:
     type: array
     items: File 
     inputBinding:
       loadContents: true
  datafilepdh:
    type:
      type: array
      items: File
      inputBinding:
        loadContents: true

outputs:
  fileprefix: string[]
  collectiondir: Directory[]

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
