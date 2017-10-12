cwlVersion: v1.0
class: ExpressionTool

requirements:
  - class: InlineJavascriptRequirement

hints:
  cwltool:LoadListingRequirement:
    loadListing: shallow_listing

inputs:
  datafilenames:
    type: File
    inputBinding:
      loadContents: true
  datafilepdh:
    type: File
    inputBinding:
      loadContents: true

outputs:
  fileprefix: string[]
  collectiondir: Directory[]

expression: |
   ${
     var names=inputs.datafilenames.contents.split('\n');
     var nblines1=names.length;
     var fileprefix=[];

     var pdhs=inputs.datafilepdh.contents.split('\n');
     var nblines2=pdhs.length;
     var collectiondir=[];
     var ssdir=[];

     for (var i = 0; i < nblines1-1; i++) {
       var ss=names[i];
       fileprefix.push(ss);
       }

     for (var i = 0; i < nblines2-1; i++) {
       var ss=pdhs[i];
       var ssdir={"class": "Directory", "location": "keep:" + ss};
       collectiondir.push(ssdir);
       }

     return {"fileprefix": fileprefix,"collectiondir":collectiondir};
     }
