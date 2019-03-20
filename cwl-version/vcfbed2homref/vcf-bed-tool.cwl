$namespaces:
  arv: "http://arvados.org/cwl#"
  cwltool: "http://commonwl.org/cwltool#"
cwlVersion: v1.0
class: CommandLineTool
requirements:
  - class: DockerRequirement
    dockerPull: vcfbed2homref
  - class: ResourceRequirement
    coresMin: 1
  - class: InitialWorkDirRequirement
    listing:
     - entry: $(inputs.ref)
       writable: True
baseCommand: bash
inputs:
  script:
    type: File
    default:
      class: File
      location: src/convert-vcf-bed-to-gvcf
    inputBinding:
      position: 1
  vcf:
    type: File
    inputBinding:
      position: 2
    secondaryFiles:
      - .tbi
  bed:
    type: File
    inputBinding:
      position: 3
  ref:
    type: File
    inputBinding:
      position: 4
      valueFrom: $(self.basename)
  out_file:
    type: string
    inputBinding:
      position: 5
outputs:
  result:
    type: File
    outputBinding:
      glob: "*.vcf.gz"
    secondaryFiles:
      - .tbi
