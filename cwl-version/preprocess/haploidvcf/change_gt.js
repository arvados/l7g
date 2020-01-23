function record() {
  var inputGT = SAMPLES[0].GT;
  if (inputGT.indexOf('/') == -1 && inputGT.indexOf('|') == -1 ) {
    SAMPLES[0].GT = inputGT + "/" + inputGT;
  } else if (CHROM == 'chrM' && inputGT.indexOf('/') != -1) {
    return inputGT.split('/')[0] == inputGT.split('/')[1];
  } else if (CHROM == 'chrM' && inputGT.indexOf('|') != -1) {
    return inputGT.split('|')[0] == inputGT.split('|')[1];
  }
}
