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

#gt_info_fn = "tmp/tile-gt-info"

ifn = sys.argv[1]
#if len(sys.argv) > 2:
#  gt_info_fn = sys.arg[2]
allele = 0
if len(sys.argv) > 2:
  allele = int(sys.argv[2])



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
cur_anchor_tilestep = ""

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

def emit(seq_map, tilepath, tilestep, _gt_list, chrom, allele):
  print_fastj = False

  tilevar = str_hex(allele, 3)

  istep = int(tilestep, 16)

  local_debug=False

  m = seq_map[ tilepath + ".00." + tilestep + "." + tilevar ]
  seq = m["seq"]
  cur_hdr = m["hdr"]

  seqa = ""

  spos = m["start_pos"]
  epos = m["end_pos"] # non-inclusive

  epos_inclusive = epos-1

  if local_debug:
    print "# start_pos:", spos, "end_pos:", epos, "seq:", seq

  gt_info_a = []

  is_spanning_tile = False
  cur_istep = istep
  seedTileLength = 1

  gt_pos = 0
  start_rel_pos = 0
  non_ref_on_end = False
  last_tilestep = tilestep
  for x in _gt_list:

    gt_pos = x["pos"]
    gt_val = x["gt"]
    gt_tilestep = x["tilestep"]
    gt_var = gt_map[gt_val[0]]

    if gt_tilestep != last_tilestep:
      seqa += seq[start_rel_pos:]
      start_rel_pos = 24

    if local_debug:
      print "# gt:", x

    m = seq_map[ tilepath + ".00." + gt_tilestep + ".000" ]
    seq = m["seq"]
    hdr = m["hdr"]

    spos = m["start_pos"]
    epos = m["end_pos"] # non-inclusive

    epos_inclusive = epos-1

    cur_hdr["endTag"] = hdr["endTag"]
    cur_hdr["endTile"] = hdr["endTile"]

    p = gt_pos - spos

    ## DEBUG
    #tm = seq_map[ tilepath + ".00." + tilestep + ".000" ]
    #anchor_spos = tm["start_pos"]
    #subseq = seq[start_rel_pos:p] + gt_var.upper()
    #print "# add pos", gt_pos, "gt", gt_val, "gtv", gt_var, "spos", spos, "start_rel_pos", start_rel_pos, "p", p, "poscheck", \
    #    anchor_spos + len(seqa) + len(subseq), len(gt_val), len(gt_var), len(subseq)



    #seqa += seq[start_rel_pos:p] + gt_var
    seqa += seq[start_rel_pos:p] + gt_var.upper()
    gt_info_a.append(p)

    #start_rel_pos = gt_pos - spos + 1
    start_rel_pos = p+1

    last_tilestep = gt_tilestep

    non_ref_on_end = False
    if (epos_inclusive - gt_pos) < 24:
      if gt_var != seq[p:p+1]:
        non_ref_on_end = True

  seqa += seq[start_rel_pos:]

  # calcluate seedTileLength from last seen tilestep and given anchor tilestep
  #
  seedTileLength = int(last_tilestep, 16) - int(tilestep, 16) + 1

  # if we have a non-ref variant on the end tag, add the remaining ref sequence
  # to this one
  #
  if non_ref_on_end and not cur_hdr["endTile"]:
    seedTileLength += 1
    m = seq_map[ tilepath + ".00." + str_hex(int(last_tilestep,16)+1,4) + ".000" ]
    seqa += m["seq"][24:]
    cur_hdr["endTag"] = hdr["endTag"]
    cur_hdr["endTile"] = hdr["endTile"]

    #print "## adding rest for ", str_hex(int(last_tilestep, 16)+1, 4), "from", last_tilestep

  if local_debug:
    print "#", md5.new(seqa).hexdigest()
    print seqa


  tileid_a = tilepath + ".00." + tilestep + "." + tilevar

  ohdr = { "tileID": tileid_a,
      "md5sum":md5.new(seqa.lower()).hexdigest(),
      "tagmask_md5sum":md5.new(seqa.lower()).hexdigest(),
      "locus": cur_hdr["locus"],
      "n": len(seqa),
      #"seedTileLength": hdr["seedTileLength"],
      "seedTileLength": seedTileLength,
      "startTile": cur_hdr["startTile"],
      "endTile": cur_hdr["endTile"],
      "startSeq": seqa[0:24],
      "endSeq": seqa[-24:],
      "startTag": cur_hdr["startTag"],
      "endTag": cur_hdr["endTag"],
      "nocallCount": nocall_count(seqa),
      "notes":[] }

  if print_fastj:
    hdr_print(ohdr)
    fold_print(seqa, 50)
    print ""
  else:
    sglf_print(ohdr, seqa)

  tile_gt_info[tileid_a] = gt_info_a

prev_chrom = ""
gt_list = []

def end_spanning_tile(seq_map, tilepath, anchor_tilestep, _gt_list, allele):
  is_spanning_tile = False

  print "####", _gt_list

  max_tilestep = _gt_list[-1]["tilestep"]

  for gt in _gt_list:

    gt_pos = gt["pos"]
    gt_val = gt_map[gt["gt"]]
    gt_tilestep = gt["tilestep"]

    if gt_tilestep != max_tilestep: continue

    m = seq_map[ tilepath + ".00." + gt_tilestep + ".000" ]
    seq = m["seq"]
    hdr = m["hdr"]
    spos = m["start_pos"]
    epos = m["end_pos"] # non-inclusive
    n = epos - spos

    rel_pos = gt_pos - spos

    #print "gt_pos", gt_pos, "gt_val", gt_val, "spos", spos, "epos", epos, "rel_pos", rel_pos, "n", n, "...", n-rel_pos

    if (n-rel_pos) >= 24: continue

    if gt_val == seq[rel_pos:rel_pos+1]:
      #print "non spanning tile:", tilepath, anchor_tilestep, cur_tilestep, gt_pos, gt_val, "==", seq[rel_pos:rel_pos+1], allele
      continue

    is_spanning_tile = True

    #print "!!spanning tile:", tilepath, anchor_tilestep, cur_tilestep, gt_pos, gt_val, "!=", seq[rel_pos:rel_pos+1], allele

  return is_spanning_tile


ifp = open(ifn)
for line in ifp:
  line = line.strip()

  # skip empty lines and comments
  #
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

  # skip indels
  #
  if (gtval[0] == 'I') or (gtval[0] == 'D'):
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
      if len(gt_list)>0:
        emit(seq_map, cur_tilepath, cur_anchor_tilestep, gt_list, prev_chrom, allele )
      cur_anchor_tilestep = tilestep
      gt_list = []

    seq_map = read_fj( fj_fn )

    cur_tilepath = tilepath
    cur_tilestep = tilestep
    cur_anchor_tilestep = tilestep
    #update = True

  # if we've skipped a tilestep, emit it.
  #
  if cur_tilestep != tilestep:
    a = int(cur_tilestep, 16)
    b = int(tilestep, 16)

    if cur_tilestep != "":
      if (b-a) > 1:
        if len(gt_list)>0:
          emit(seq_map, cur_tilepath, cur_anchor_tilestep, gt_list, prev_chrom, allele )
        gt_list = []
        cur_anchor_tilestep = tilestep
      elif (len(gt_list)>0) and (not end_spanning_tile(seq_map, cur_tilepath, cur_anchor_tilestep, gt_list, allele)):
        emit(seq_map, cur_tilepath, cur_anchor_tilestep, gt_list, prev_chrom, allele )
        gt_list = []
        cur_anchor_tilestep = tilestep
      elif len(gt_list)==0:
        cur_anchor_tilestep = tilestep

  if allele < len(gtval):

    gt_list.append( { "pos": pos0ref, "gt" : gtval[allele], "tilestep": tilestep } )

    ## DEBUG
    #print "#gt_list++", pos0ref, gtval[allele], tilepath, tilestep, cur_tilestep, cur_anchor_tilestep, gt_list

  cur_tilestep = tilestep

  prev_chrom = chrom

  ####DEBUG
  #if tilepath == "0001": sys.exit()

if len(gt_list)>0:
  print "#fin"
  emit(seq_map, cur_tilepath, cur_anchor_tilestep, gt_list, prev_chrom, allele)

ifp.close()

#fp = open(gt_info_fn, "w")
#for key in tile_gt_info:
#  a = tile_gt_info[key]
#  fp.write(str(key))
#  for v in a:
#    fp.write(" " + str(v))
#  fp.write("\n")
#fp.close()
