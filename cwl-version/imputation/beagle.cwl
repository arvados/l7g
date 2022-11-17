cwlVersion: v1.1
class: CommandLineTool
requirements:
  ShellCommandRequirement: {}
hints:
  DockerRequirement:
    dockerPull: beagle5.4
  ResourceRequirement:
    coresMin: 2
    ramMin: 10000
inputs:
  sample: string
  chr: string
  ref: File
  map: File
  vcf:
    type: File
    secondaryFiles: [.tbi]
outputs:
  rawimputedvcf:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles: [.tbi]
baseCommand: [bcftools, view]
arguments:
  - $(inputs.vcf)
  - prefix: "--regions"
    valueFrom: $(inputs.chr)
  - "-Oz"
  - prefix: "-o"
    valueFrom: $(inputs.sample)_$(inputs.chr).vcf.gz
  - shellQuote: false
    valueFrom: "&&"
  - "java"
  - -Xms$(runtime.ram)m
  - prefix: "-jar"
    valueFrom: "/beagle.05May22.33a.jar"
  - prefix: "ref="
    separate: false
    valueFrom: $(inputs.ref)
  - prefix: "map="
    separate: false
    valueFrom: $(inputs.map)
  - prefix: "gt="
    separate: false
    valueFrom: $(inputs.sample)_$(inputs.chr).vcf.gz
  - prefix: "out="
    separate: false
    valueFrom: $(inputs.sample)_rawimputed_$(inputs.chr)
  - prefix: "nthreads="
    separate: false
    valueFrom: $(runtime.cores)
  - shellQuote: false
    valueFrom: "&&"
  - "tabix"
  - $(inputs.sample)_rawimputed_$(inputs.chr).vcf.gz
  - shellQuote: false
    valueFrom: "&&"
  - "rm"
  - $(inputs.sample)_$(inputs.chr).vcf.gz
