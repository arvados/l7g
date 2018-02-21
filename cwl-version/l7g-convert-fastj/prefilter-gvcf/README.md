CWL Lightning VCF Filter
===

Sometimes the gVCF (VCF v4.1 or above) has data
that is incompatible with how we want to process
gVCF.
This CWL pipeline will preformat the gVCF files
into a format that tools downstream expect.
This will also `bgzip` and index the gVCF
files suitable for tabix usage.
