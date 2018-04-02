class: ExpressionTool
cwlVersion: v1.0
requirements:
  InlineJavascriptRequirement: {}

inputs:
  idir: Directory
  nBatch: string
outputs:
  batchDir: Directory[]
  batchName: string[]
expression: |
  ${
    var batched_dir = [];
    var n = inputs.idir.listing.length;
    var m = ( (typeof inputs.nBatch === "undefined") ? 10 : parseInt(inputs.nBatch) );
    if (m<1) { m = 1; }
    var m = Math.ceil( n / m );

    var dat = { "batchDir" : [], "batchName": [] };

    var batch_dir = {};
    for (var idx=0; idx<inputs.idir.listing.length; idx++) {

      if ((m==1) || ((idx % m) == 0)) {
        if (idx>0) {
          var bn = dat.batchName.length;
          dat.batchDir.push(batch_dir);
          dat.batchName.push( "batch" + bn.toString() );
        }
        batch_dir = { "class" : "Directory", "basename": ".", "listing" : [] };
      }
      var ele = inputs.idir.listing[idx];
      batch_dir.listing.push(ele);
    }
    if ((typeof batch_dir.listing !== "undefined") &&
        (batch_dir.listing.length > 0)) {
      var bn = dat.batchName.length;
      dat.batchDir.push(batch_dir);
      dat.batchName.push( "batch" + bn.toString() );
    }
    return dat;
  }
