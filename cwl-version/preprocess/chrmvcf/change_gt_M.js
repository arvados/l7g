function record() {
  if (CHROM == 'M') {
    var inputGT = SAMPLES[0].GT;
    if (inputGT.indexOf('/') == -1 && inputGT.indexOf('|') == -1 ) {
      SAMPLES[0].GT = inputGT + "/" + inputGT;
    } else if (inputGT.indexOf('/') != -1 && inputGT.split('/')[0] != inputGT.split('/')[1]) {
      return false;
    } else if (inputGT.indexOf('|') != -1 && inputGT.split('|')[0] != inputGT.split('|')[1]) {
      return false;
    }
  }
}
