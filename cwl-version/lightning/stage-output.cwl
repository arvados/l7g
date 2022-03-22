cwlVersion: v1.2
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
hints:
  LoadListingRequirement:
    loadListing: shallow_listing
inputs:
  libname: string
  npyfiles: File[]
  onehotnpyfiles: File[]
  bed: File?
  annotatedvcf: File?
  summary: File?
outputs:
  stagednpydir: Directory
  stagedonehotnpydir: Directory
  stagedannotationdir: Directory?
expression: |
  ${
    var stagednpydir = {"class": "Directory",
                        "basename": "library_full",
                        "listing": inputs.npyfiles};
    var stagedonehotnpydir = {"class": "Directory",
                              "basename": "library_filtered",
                              "listing": inputs.onehotnpyfiles};
    var annotationlist = [];
    if (inputs.bed != "null") {
      annotationlist.push(inputs.bed);
    }
    if (inputs.annotatedvcf != "null") {
      annotationlist.push(inputs.annotatedvcf);
    }
    if (inputs.summary != "null") {
      annotationlist.push(inputs.summary);
    }
    if (annotationlist != []) {
      var stagedannotationdir = {"class": "Directory",
                                 "basename": inputs.libname+"_annotation",
                                 "listing": annotationlist};
    }
    return {"stagednpydir": stagednpydir, "stagedonehotnpydir": stagedonehotnpydir, "stagedannotationdir": stagedannotationdir};
  }
