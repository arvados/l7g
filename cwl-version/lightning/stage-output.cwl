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
  annotatedvcf: File?
  summary: File?
outputs:
  stagednpydir: Directory
  stagedonehotnpydir: Directory
  stagedannotationdir: Directory?
expression: |
  ${
    var stagednpydir = {"class": "Directory",
                        "basename": inputs.libname+"_npy",
                        "listing": inputs.npyfiles};
    var stagedonehotnpydir = {"class": "Directory",
                              "basename": inputs.libname+"_onehotnpy",
                              "listing": inputs.onehotnpyfiles};
    if (inputs.annotatedvcf != "null") {
      var annotationlist = [inputs.annotatedvcf];
      if (inputs.summary != "null") {
        annotationlist.push(inputs.summary);
      }
      var stagedannotationdir = {"class": "Directory",
                                 "basename": inputs.libname+"_annotation",
                                 "listing": annotationlist};
    } else {
      var stagedannotationdir = "null";
    }
    return {"stagednpydir": stagednpydir, "stagedonehotnpydir": stagedonehotnpydir, "stagedannotationdir": stagedannotationdir};
  }
