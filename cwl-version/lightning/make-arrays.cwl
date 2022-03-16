cwlVersion: v1.2
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
inputs:
  matchgenome_array:
    type: string[]
  libdir_array:
    type: Directory[]
  genomeversion_array:
    type: string[]
  regions_nestedarray:
    type:
      type: array
      items:
        type: array
        items: [File, "null"]
  threads_array:
    type: int[]
  mergeoutput_array:
    type: string[]
  expandregions_array:
    type: int[]
outputs:
  full_matchgenome_array:
    type: string[]
  full_libdir_array:
    type: Directory[]
  full_genomeversion_array:
    type: string[]
  full_regions_array:
    type:
      type: array
      items: [File, "null"]
  full_threads_array:
    type: int[]
  full_mergeoutput_array:
    type: string[]
  full_expandregions_array:
    type: int[]
  full_libname_array:
    type: string[]
expression: |
  ${
    var full_matchgenome_array = [];
    var full_libdir_array = [];
    var full_genomeversion_array = [];
    var full_regions_array = [];
    var full_threads_array = [];
    var full_mergeoutput_array = [];
    var full_expandregions_array = [];
    var full_libname_array = [];
    for (var i = 0; i < inputs.matchgenome_array.length; i++) {
      for (var j = 0; j < inputs.libdir_array.length; j++) {
        for (var k = 0; k < inputs.regions_nestedarray[j].length; k++) {
          full_matchgenome_array.push(inputs.matchgenome_array[i]);
          full_libdir_array.push(inputs.libdir_array[j]);
          full_genomeversion_array.push(inputs.genomeversion_array[j]);
          full_regions_array.push(inputs.regions_nestedarray[j][k]);
          full_threads_array.push(inputs.threads_array[k]);
          full_mergeoutput_array.push(inputs.mergeoutput_array[k]);
          full_expandregions_array.push(inputs.expandregions_array[k]);
          var libname = inputs.genomeversion_array[j]+inputs.matchgenome_array[i]+"_library";
          full_libname_array.push(libname);
        }
      }
    }
    return {"full_matchgenome_array": full_matchgenome_array, 
            "full_libdir_array": full_libdir_array, "full_genomeversion_array": full_genomeversion_array,
            "full_regions_array": full_regions_array, "full_threads_array": full_threads_array, "full_mergeoutput_array": full_mergeoutput_array, "full_expandregions_array": full_expandregions_array,
            "full_libname_array": full_libname_array};
  }
