#!/usr/bin/python
#
# convert the space delimited file
# created by `gt-23andMeConvert.py`
# to fastj or sglf
#

import os
import sys
import gzip
import json
import md5


if len(sys.argv) < 2:
  print "provide input ssv file"
  sys.exit(0)

gt_info_fn = "tmp/tile-gt-info"

ifn = sys.argv[1]
if len(sys.argv) > 2:
  gt_info_fn = sys.arg[2]


fjdir = "/data-sdd/data/fastj/hg19"

fj_base = "/data-sdd/data/fastj/hg19"
fj_fn = "/data-sdd/data/fastj/hg19/0000.fj.gz"

#hg19_band_base = "/data-sdd/scripts/cgf3/hg19.band"
#hg19_band0 = []
#hg19_band1 = []

# save information about where the variant appears in the tile
# for later output to a file.
#
tile_gt_info = {}

cur_tilepath = ""
cur_tilestep = ""

def nocall_count(seq):
  count = 0
  for x in seq:
    if x == 'n' or x == 'N':
      count+=1
  return count

def hdr_print(hdr):
  ord_key = ["tileID", "md5sum", "tagmask_md5sum", "locus", "n", "seedTileLength",
      "startTile","endTile","startSeq","endSeq","startTag","endTag",
      "nocallCount","notes"]

  count = 0
  o_str = "> {"
  for key in ord_key:
    if count>0: o_str+=","
    o_str += '"' + key + '":' + json.dumps(hdr[key])
    count+=1
  o_str += "}"
  print o_str


def sglf_print(hdr, seq):
  print hdr["tileID"] + "+" + str(hdr["seedTileLength"]) + "," + hdr["md5sum"] + "," + seq

def fold_print(s, w):
  for x in range(0, len(s), w):
    e = x+w
    if e>len(s):
      e = len(s)
    print s[x:e]

def read_band(fn, band0, band1):
  ifp = open(fn)
  count = 0
  for line in ifp:
    t = line.strip()
    sval = t[1:-1].strip().split(" ")
    if count == 0:
      for v in sval:
        band0.append(int(v))
    else:
      for v in sval:
        band1.append(int(v))
    count += 1
    if count >= 2: break

  ifp.close()

#x = []
#y = []
#read_band(hg19_band_base + "/" + "hg19.0db", x, y)

def read_fj(fn):
  """read in the gzip'd fastj file and store
  the tile sequences in a map whose key is
  the tileID in 4,2,4.3 format (e.g. 0000.00.013c.001)
  """

  seq_lookup = {}

  fj_fp = gzip.open(fn)

  cur_tileid = ""
  cur_seq = ""

  build_start_pos = -1
  build_end_pos = -1

  hdr_json = {}

  for line in fj_fp:
    line = line.strip()
    if len(line)==0: continue

    if line[0]=='>':
      if len(cur_tileid)>0:

        # END IS NON INCLUSIVE (PYTHON DEFAULT)
        seq_lookup[cur_tileid] = { "seq":cur_seq, "start_pos":build_start_pos, "end_pos":build_end_pos, "hdr":hdr_json}

      hdr_json = json.loads(line[1:])
      cur_tileid = hdr_json["tileID"]
      cur_seq = ""

      build_info = hdr_json["locus"][0]["build"].split(" ")

      build_start_pos = int(build_info[2])
      build_end_pos = int(build_info[3])

      continue

    cur_seq += line


  seq_lookup[cur_tileid] = { "seq":cur_seq, "start_pos":build_start_pos, "end_pos":build_end_pos, "hdr": hdr_json}

  fj_fp.close()

  # debug
  #for key in seq_lookup:
  #  print key, md5.new(seq_lookup[key]["seq"]).hexdigest(), seq_lookup[key]

  return seq_lookup

# debug
#sm = read_fj(fj_fn)
#sys.exit(0)

gt_map = { "A" : "a", "a" :"a",
    "C":"c", "c":"c",
    "G":"g", "g":"g",
    "T":"t", "t":"t",
    "-":"n", "N":"n", "n":"n",
    "D":"", "I":"", "d":"", "i":"" }

seq_map = {}

def str_hex(v, l):
  s = hex(v)[2:]
  n = "0"*(l-len(s))
  return n + s

def emit(seq_map, tilepath, tilestep, gt_list, chrom, ovariant):
  print_fastj = False

  istep = int(tilestep, 16)

  local_debug=False

  m = seq_map[ tilepath + ".00." + tilestep + ".000" ]
  seq = m["seq"]
  hdr = m["hdr"]

  if local_debug:
    print "# tile:", tilepath + ".00." + tilestep + ".000"

  seqa = ""

  spos = m["start_pos"]
  epos = m["end_pos"] # non-inclusive

  epos_inclusive = epos-1

  if local_debug:
    print "# start_pos:", spos, "end_pos:", epos, "seq:", seq

  gt_info_a = []

  spanning_tile = False
  cur_istep = istep
  seedTileLength =1

  gt_pos = 0
  prev_rel_pos = 0
  for x in gt_list:

    if local_debug:
      print "# gt:", x

    if (epos_inclusive - gt_pos) < 24:
      spanning_tile = True

    if (gt_pos >= epos):
      seqa += seq[prev_rel_pos:]
      if allele>0:
        seqb += seq[prev_rel_pos:]
      cur_istep += 1 # hg19, seedTileLength=1

      m = seq_map[ tilepath + ".00." + str_hex(cur_istep, 4)+ ".000" ]
      seq = m["seq"]
      hdr["endTag"] = m["endTag"]
      hdr["endTile"] = m["endTile"]
      seedTileLength += 1



    gt_pos = x["pos"]
    gt_val = x["gt"]
    variant = [ gt_map[gt_val[0]], "" ]

    if variant[0] != "":
      p = x["pos"] - spos
      seq += seq[prev_rel_pos:p] + variant[0]
      gt_info_a.append(p)

    prev_rel_pos = gt_pos - spos + 1


  seqa += seq[prev_rel_pos:]

  if local_debug:
    print "#", md5.new(seqa).hexdigest()
    print seqa

    if allele>1:
      print "#", md5.new(seqb).hexdigest()
      print seqb


  tileid_a = tilepath + ".00." + tilestep + "." + ovariant

  ohdr = { "tileID": tileid_a,
      "md5sum":md5.new(seqa).hexdigest(),
      "tagmask_md5sum":md5.new(seqa).hexdigest(),
      "locus": hdr["locus"],
      "n": len(seqa),
      #"seedTileLength": hdr["seedTileLength"],
      "seedTileLength": seedTileLength,
      "startTile": hdr["startTile"],
      "endTile": hdr["endTile"],
      "startSeq": seqa[0:24],
      "endSeq": seqa[-24:],
      "startTag": hdr["startTag"],
      "endTag": hdr["endTag"],
      "nocallCount": nocall_count(seqa),
      "notes":[] }

  if print_fastj:
    hdr_print(ohdr)
    fold_print(seqa, 50)
    print ""
  else:
    sglf_print(ohdr, seqa)

  #print "### adding", tileid_a, gt_info_a
  tile_gt_info[tileid_a] = gt_info_a
  if len(gt_val)>1:
    #print "### adding", tileid_b, gt_info_b
    tile_gt_info[tileid_b] = gt_info_b

def emit_old(seq_map, tilepath, tilestep, gt_list, chrom):
  print_fastj = False

  istep = int(tilestep, 16)

  allele = 2
  if chrom == "chrM" or chrom == "chrX" or chrom == "chrY":
    allele = 1

  local_debug=False

  m = seq_map[ tilepath + ".00." + tilestep + ".000" ]
  seq = m["seq"]
  hdr = m["hdr"]

  if local_debug:
    print "# tile:", tilepath + ".00." + tilestep + ".000"

  seqa = ""
  seqb = ""

  spos = m["start_pos"]
  epos = m["end_pos"] # non-inclusive

  epos_inclusive = epos-1

  if local_debug:
    print "# start_pos:", spos, "end_pos:", epos, "seq:", seq

  gt_info_a = []
  gt_info_b = []

  spanning_tile = False
  cur_istep = istep
  seedTileLength =1

  prev_rel_pos = 0
  for x in gt_list:

    if local_debug:
      print "# gt:", x

    if (epos_inclusive - gt_pos) < 24:
      spanning_tile = True

    if (gt_pos >= epos):
      seqa += seq[prev_rel_pos:]
      if allele>0:
        seqb += seq[prev_rel_pos:]
      cur_istep += 1 # hg19, seedTileLength=1

      m = seq_map[ tilepath + ".00." + str_hex(cur_istep, 4)+ ".000" ]
      seq = m["seq"]
      hdr["endTag"] = m["endTag"]
      hdr["endTile"] = m["endTile"]
      seedTileLength += 1



    gt_pos = x["pos"]
    gt_val = x["gt"]
    variant = [ gt_map[gt_val[0]], "" ]

    if len(gt_val) > 1:
      variant[1] = gt_map[gt_val[1]]

    if variant[0] != "":
      p = x["pos"] - spos
      seqa += seq[prev_rel_pos:p] + variant[0]
      gt_info_a.append(p)


    if allele>1:
      if variant[1] != "":
        p = x["pos"] - spos
        seqb += seq[prev_rel_pos:p] + variant[1]
        gt_info_b.append(p)

    prev_rel_pos = gt_pos - spos + 1


  seqa += seq[prev_rel_pos:]

  if allele>0:
    seqb += seq[prev_rel_pos:]

  if local_debug:
    print "#", md5.new(seqa).hexdigest()
    print seqa

    if allele>1:
      print "#", md5.new(seqb).hexdigest()
      print seqb


  tileid_a = tilepath + ".00." + tilestep + ".000"

  ohdr = { "tileID": tilepath + ".00." + tilestep + ".000",
      "md5sum":md5.new(seqa).hexdigest(),
      "tagmask_md5sum":md5.new(seqa).hexdigest(),
      "locus": hdr["locus"],
      "n": len(seqa),
      #"seedTileLength": hdr["seedTileLength"],
      "seedTileLength": seedTileLength,
      "startTile": hdr["startTile"],
      "endTile": hdr["endTile"],
      "startSeq": seqa[0:24],
      "endSeq": seqa[-24:],
      "startTag": hdr["startTag"],
      "endTag": hdr["endTag"],
      "nocallCount": nocall_count(seqa),
      "notes":[] }

  if print_fastj:
    hdr_print(ohdr)
    fold_print(seqa, 50)
    print ""
  else:
    sglf_print(ohdr, seqa)

  tileid_b = tilepath + ".00." + tilestep + ".001"

  if allele>1:

    ohdr = { "tileID": tilepath + ".00." + tilestep + ".001",
        "md5sum":md5.new(seqb).hexdigest(),
        "tagmask_md5sum":md5.new(seqb).hexdigest(),
        "locus": hdr["locus"],
        "n": len(seqb),
        #"seedTileLength": hdr["seedTileLength"],
        "seedTileLength": seedTileLength,
        "startTile": hdr["startTile"],
        "endTile": hdr["endTile"],
        "startSeq": seqa[0:24],
        "endSeq": seqa[-24:],
        "startTag": hdr["startTag"],
        "endTag": hdr["endTag"],
        "nocallCount": nocall_count(seqb),
        "notes":[] }

    if print_fastj:
      hdr_print(ohdr)
      fold_print(seqb, 50)
      print ""
    else:
      sglf_print(ohdr, seqb)

  #print "### adding", tileid_a, gt_info_a
  tile_gt_info[tileid_a] = gt_info_a
  if len(gt_val)>1:
    #print "### adding", tileid_b, gt_info_b
    tile_gt_info[tileid_b] = gt_info_b

prev_chrom = ""
gt_list_a = []
gt_list_b = []

ifp = open(ifn)
for line in ifp:
  line = line.strip()
  if len(line)==0: continue
  if line[0] == '#': continue

  fields = line.split(" ")

  chrom = fields[0]
  pos0ref = int(fields[1])
  tilepath = fields[2]
  tilestep = fields[3]
  gtval = fields[4]

  # skip nocalls
  #
  if gtval == "--" or gtval == "-":
    continue

  #print "# got:", fields

  if prev_chrom == "":
    prev_chrom = chrom

  if tilepath != cur_tilepath:
    fj_fn = fj_base + "/" + tilepath + ".fj.gz"

    #print "# reading", fj_fn

    # we're not on the first tile path so we have a previous entry
    # we need to emit.  emit it.
    #
    if cur_tilepath != "":
      emit(seq_map, cur_tilepath, cur_tilestep, gt_list_a, prev_chrom, "000" )
      emit(seq_map, cur_tilepath, cur_tilestep, gt_list_b, prev_chrom, "001" )

    seq_map = read_fj( fj_fn )

    cur_tilepath = tilepath
    cur_tilestep = tilestep
    #update = True

  # if we've skipped a tilestep, emit it.
  #
  if cur_tilestep != tilestep:
    if cur_tilestep != "":
      emit(seq_map, cur_tilepath, cur_tilestep, gt_list_a, prev_chrom, "000" )
      emit(seq_map, cur_tilepath, cur_tilestep, gt_list_b, prev_chrom, "001" )
    gt_list_a = []
    gt_list_b = []

  gt_list_a.append( { "pos": pos0ref, "gt" : gtval[0] } )
  if len(gtval)>1:
    gt_list_b.append( { "pos": pos0ref, "gt" : gtval[1] } )
  cur_tilestep = tilestep

  prev_chrom = chrom

  ####DEBUG
  if tilepath == "0001":
    sys.exit()

if len(gt_list)>0:
  emit(seq_map, cur_tilepath, cur_tilestep, gt_list_a, prev_chrom, "000")
  emit(seq_map, cur_tilepath, cur_tilestep, gt_list_a, prev_chrom, "001")

ifp.close()

fp = open(gt_info_fn, "w")
for key in tile_gt_info:
  a = tile_gt_info[key]
  fp.write(str(key))
  for v in a:
    fp.write(" " + str(v))
  fp.write("\n")
fp.close()
