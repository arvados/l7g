cwlVersion: v1.1
class: Workflow
requirements:
  ScatterFeatureRequirement: {}

inputs:
  sample: string
  chrs: string[]
  vcfdir: Directory
  snpeffdatadir: Directory
  genomeversion: string
  dbsnp:
    type: File
    secondaryFiles: [.csi]
  gnomaddir: Directory

outputs:
  annotatedvcf:
    type: File
    secondaryFiles: [.tbi]
    outputSource: bcftools-concat/vcf
  summary:
    type: File
    outputSource: totalcounts/summary

steps:
  getfiles:
    run: getfiles.cwl
    in:
      sample: sample
      chrs: chrs
      vcfdir: vcfdir
      gnomaddir: gnomaddir
    out: [samples, vcfs, gnomads]

  preprocess:
    run: preprocess.cwl
    scatter: [sample, vcf]
    scatterMethod: dotproduct
    in:
      sample: getfiles/samples
      vcf: getfiles/vcfs
    out: [trimmedvcf]

  snpeff-bcftools-annotate:
    run: snpeff-bcftools-annotate.cwl
    scatter: [sample, vcf, gnomad]
    scatterMethod: dotproduct
    in:
      vcf: preprocess/trimmedvcf
      sample: getfiles/samples
      snpeffdatadir: snpeffdatadir
      genomeversion: genomeversion
      dbsnp: dbsnp
      gnomad: getfiles/gnomads
    out: [annotatedvcf]

  bcftools-concat:
    run: bcftools-concat.cwl
    in:
      sample: sample
      vcfs: snpeff-bcftools-annotate/annotatedvcf
    out: [vcf]

  getcount:
    run: getcount.cwl
    scatter: [sample, vcf]
    scatterMethod: dotproduct
    in:
      sample: getfiles/samples
      vcf: snpeff-bcftools-annotate/annotatedvcf
    out: [count]

  totalcounts:
    run: totalcounts.cwl
    in:
      sample: sample
      counts: getcount/count
    out: [summary]
