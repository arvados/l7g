cwlVersion: v1.2
class: Workflow
requirements:
  InlineJavascriptRequirement: {}
  SubworkflowFeatureRequirement: {}
  MultipleInputFeatureRequirement: {}

inputs:
  matchgenome: string
  libdir: Directory
  regions: File?
  threads: int
  mergeoutput: string
  expandregions: int
  phenotypesdir: Directory
  libname: string
  chrs: string[]
  snpeffdatadir: Directory
  genomeversion: string
  dbsnp:
    type: File
    secondaryFiles: [.csi]
  gnomaddir: Directory
  readmeinfo: string[]

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
  readme:
    type: File
    outputSource: genreadme/readme

steps:
  lightning-tiling-stats:
    run: lightning-tiling-stats.cwl
    when: $(inputs.regions == null)
    in:
      libdir: libdir
    out: [bed]

  lightning-slice-numpy:
    run: lightning-slice-numpy.cwl
    in:
      matchgenome: matchgenome
      libdir: libdir
      regions: regions
      threads: threads
      mergeoutput: mergeoutput
      expandregions: expandregions
    out: [outdir, npys, samplescsv, chunktagoffsetcsv]

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
    out: [outdir, npys, samplescsv]

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
        source: [lightning-slice-numpy/npys, lightning-slice-numpy/samplescsv, lightning-slice-numpy/chunktagoffsetcsv]
        linkMerge: merge_flattened
      onehotnpyfiles:
        source: [lightning-slice-numpy-onehot/npys, lightning-slice-numpy-onehot/samplescsv]
        linkMerge: merge_flattened
      bed: lightning-tiling-stats/bed
      annotatedvcf: annotate-wf/annotatedvcf
      summary: annotate-wf/summary
    out: [stagednpydir, stagedonehotnpydir, stagedannotationdir]

  genreadme:
    run: genreadme.cwl
    in:
      samplescsv: lightning-slice-numpy/samplescsv
      readmeinfo: readmeinfo
    out: [readme]
