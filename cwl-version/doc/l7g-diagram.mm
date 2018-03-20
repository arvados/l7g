graph TD

subgraph REF_input
hg19[hg19.fa.gz]
hg38[hg38.fa.gz]
cytoband[cytoband.hg19.custom.txt]
wgEncode[wgEncodeCrgMapabilityAlign24mer.bw]
end

hg19 -.-> tagset_wf
hg38 -.-> tile_assembly_hg38_wf
cytoband -.-> tagset_wf
wgEncode -.-> tagset_wf

hg19 -.-> tile_assembly_hg19_wf
cytoband -.-> tile_assembly_hg19_wf

subgraph tagset_workflow
tagset_wf(tagset workflow) --> tagset_data[tagset.fa.gz]
end
tagset_data -.-> tile_assembly_hg19_wf

subgraph tile_assembly_hg19_workflow
tile_assembly_hg19_wf(tile assembly hg19) --> tile_assembly_hg19_data[assembly.00.hg19.fw.gz]
end
tile_assembly_hg19_data -.-> tile_assembly_hg38_wf

subgraph tile_assembly_hg38_workflow
tile_assembly_hg38_wf(tile assembly hg38) --> tile_assembly_hg38_data[assembly.00.hg38.fw.gz]
end

subgraph GFF_input
hupgp_gff0[HuPGP GFF 0]
hupgp_gff1[HuPGP GFF 1]
hupgp_gff_dotdotdot[HuPGP GFF ...]
hupgp_gff_n[HuPGP GFF n-1 ]
end

hupgp_gff0 -.-> hupgp_fastj_data0
hupgp_gff1 -.-> hupgp_fastj_data1
hupgp_gff_dotdotdot -.-> hupgp_fastj_data_dotdotdot
hupgp_gff_n -.-> hupgp_fastj_data_n

hg19 -.-> hupgp_fastj_wf0
hg19 -.-> hupgp_fastj_wf1
hg19 -.-> hupgp_fastj_wf_dotdotdot
hg19 -.-> hupgp_fastj_wf_n

tagset_data -.-> hupgp_fastj_wf0
tagset_data -.-> hupgp_fastj_wf1
tagset_data -.-> hupgp_fastj_wf_dotdotdot
tagset_data -.-> hupgp_fastj_wf_n

subgraph HuPGP_FastJ
hupgp_fastj_wf0(HuPGP FastJ 0) --> hupgp_fastj_data0[HuPGP FastJ 0]
hupgp_fastj_wf1(HuPGP FastJ 1) --> hupgp_fastj_data1[HuPGP FastJ 1]
hupgp_fastj_wf_dotdotdot(HuPGP Fastj ...) --> hupgp_fastj_data_dotdotdot[HuPGP FastJ ...]
hupgp_fastj_wf_n(HuPGP FastJ n-1) --> hupgp_fastj_data_n[HuPGP FastJ n-1]

hupgp_fastj_data0 -.-> hupgp_fastj_gathered_data[HuPGP FastJ All]
hupgp_fastj_data1 -.-> hupgp_fastj_gathered_data
hupgp_fastj_data_dotdotdot -.-> hupgp_fastj_gathered_data
hupgp_fastj_data_n -.-> hupgp_fastj_gathered_data
end

