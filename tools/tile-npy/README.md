Lightning Numpy Tile Arrays
===

This directory holds programs to help with the creation of lightning tile arrays.

band-to-matrix-npy
---

A program to convert band information into a Lightning tile numpy array.
This includes low quality positions and high quality positions.

### Quick start

```
$ make band-to-matrix-npy
$ cat dataset0-003.band dataset1-003.band dataset2-003.band | ./band-to-matrix-npy 003 003.npy
```

npy-vec-to-hiq-1hot
---

A program to filter only high quality tiles and to create a corresponding 1hot
numpy matrix of the high quality tiles.
"Knots" will be removed if any tile in the knot has a low quality tile in them.

Note that this requires a potentially high memory machine as the whole matrix is
stored in memory before writing to disk.

### Quick start

```
$ make npy-vec-to-hiq-1hot
$ ./npy-vec-to-hiq-1hot names.npy inp_vec_npy_dir/ out_dir/
```

Where `names.npy` holds a numpy array of the names of the datasets, `inp_vec_npy_dir`
is the location of the original Lightning tile arrays (produced by `band-to-matrix-npy`, say)
and `out_dir` holds the location of the directory to put the resulting high quaility lightning
tile numpy arrays.

This will create:

* `hiq`
* `hiq-info`
* `hiq-1hot`
* `hiq-1hot-info`

npy-consolidate
---

A program to consolodate the Lightning tile numpy arrays into a single matrix.
`band-to-matrix-npy` creates a Lightning tile numpy array on a tilepath by tilepath
basis.
This program consolodates them into a single matrix.

### Quick start

```
$ make npy-consolidate
$ ./npy-consolidate inp_vec_npy_dir/000 inp_vec_npy_dir/001 ... inp_vec_npy_dir/035e out_dir/all
```

This will create an `all` numpy matrix as well as an `all-info` numpy matrix that encodes the
Lightning tile path and Lightning tile step of each column in the `all` numpy matrix.
