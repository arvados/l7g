Tile Assembly (hg19)
===

This is the Common Workflow Language (CWL) pipeline
to create the tile assembly for hg19.
The tile assembly maps tile boundaries to a reference
genome, in this case, hg19.

Care needs to be taken when trying to create tile
assemblies for other references and this is why
this pipeline only deals with a single tile assembly
(hg19).

CWL Pipeline Submission
---

On Arvados:

```
arvados-cwl-runner --disable-reuse --local l7g-tagset.cwl l7g-tagset.yml
```

Local Run
---

```
./build-tagset.sh \
  /path/to/cytoband-file \
  /path/to/wgEncodeCrgMapabilityAlign24mer.bw \
  /path/to/hg19.fa.gz
```

This wil create a temporary `stage` directory that will be deleted
after completion.

Four files are created:

* `assembly.00.hg19.fw.gz` - The tile assembly file compressed with `bgzip`
* `assembly.00.hg19.fw.gz.gzi` - The index to the `.gz` file above
* `assembly.00.hg19.fw.gz.fwi` - The index to the fixed width
* `assembly.00.hg19.fw.fwi` - The index to the fixed width (identical to `assembly.00.hg19.fw.gz.fwi`)

