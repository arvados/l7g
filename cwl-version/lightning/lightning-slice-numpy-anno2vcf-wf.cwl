cwlVersion: v1.2
class: Workflow
requirements:
  InlineJavascriptRequirement: {}
  SubworkflowFeatureRequirement: {}
  MultipleInputFeatureRequirement: {}

inputs:
  matchgenome:
    type: string
  libdir:
    type: Directory
  regions:
    type: File?
  threads:
    type: int
  mergeoutput:
    type: string
  expandregions:
    type: int
  phenotypesdir:
    type: Directory
  libname: string
  chrs: string[]
  snpeffdatadir: Directory
  genomeversion: string
  dbsnp:
    type: File
    secondaryFiles: [.csi]
  gnomaddir: Directory

outputs:
  stagednpydir:
    type: Directory
    outputSource: stage-output/stagednpydir
  stagedonehotnpydir:
    type: Directory
    outputSource: stage-output/stagedonehotnpydir
  stagedannotationdir:
    type: Directory?
    outputSource: stage-output/stagedannotationdir

steps:
  lightning-slice-numpy:
    run: lightning-slice-numpy.cwl
    in:
      matchgenome: matchgenome
      libdir: libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
    out: [outdir, npys, csv]

  lightning-slice-numpy-onehot:
    run: lightning-slice-numpy-onehot.cwl
    in:
      matchgenome: matchgenome
      libdir: libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
      phenotypesdir: phenotypesdir
    out: [outdir, onehotcolumnsnpy, onehotnpy, csv]

  lightning-anno2vcf-onehot:
    run: lightning-anno2vcf.cwl
    when: $(inputs.regions == null)
    in:
      annodir: lightning-slice-numpy-onehot/outdir
      regions: regions
    out: [vcfdir]

  annotate-wf:
    run: ../annotation/annotate-wf.cwl
    when: $(inputs.vcfdir != null)
    in:
      sample: libname
      chrs: chrs
      vcfdir: lightning-anno2vcf-onehot/vcfdir
      snpeffdatadir: snpeffdatadir
      genomeversion: genomeversion
      dbsnp: dbsnp
      gnomaddir: gnomaddir
    out: [annotatedvcf, summary]

  stage-output:
    run: stage-output.cwl
    in:
      libname: libname
      npyfiles:
        source: [lightning-slice-numpy/npys, lightning-slice-numpy/csv]
        linkMerge: merge_flattened
      onehotnpyfiles:
        source: [lightning-slice-numpy-onehot/onehotcolumnsnpy, lightning-slice-numpy-onehot/onehotnpy, lightning-slice-numpy-onehot/csv]
        linkMerge: merge_flattened
      annotatedvcf: annotate-wf/annotatedvcf
      summary: annotate-wf/summary
    out: [stagednpydir, stagedonehotnpydir, stagedannotationdir]
