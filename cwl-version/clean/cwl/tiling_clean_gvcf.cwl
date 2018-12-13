$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: Workflow
label: Resolve duplicate/overlapping calls in gVCFs
doc: |
   Parses all gVCFs and cleans them returning a clean set
requirements:
  - class: DockerRequirement
    dockerPull: javatools
  - class: ResourceRequirement
    coresMin: 2
    coresMax: 2
  - class: ScatterFeatureRequirement
  - class: InlineJavascriptRequirement
  - class: SubworkflowFeatureRequirement
hints:
  arv:RuntimeConstraints:
    keep_cache: 4096
inputs:
  refdirectory:
    type: Directory
    label: Location in Arvados Keep of gVCFs to clean
  bashscript:
    type: File
    label: Master script to control cleaning
  cleanvcf:
    type: File
    label: Compiled code that cleans gVCFs

outputs:
  out1:
    type: Directory[]
    outputSource: step2/out1
    label: Output directory of clean gVCFs

steps:
  step1:
    run: getdirs.cwl
    in:
      refdirectory: refdirectory
    out: [out1,out2]

  step2:
    scatter: [gvcfDir,gvcfPrefix]
    scatterMethod: dotproduct
    in:
      bashscript: bashscript
      gvcfDir: step1/out1
      gvcfPrefix: step1/out2
      cleanvcf: cleanvcf
    run: cleangvcf.cwl
    out: [out1]
