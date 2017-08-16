#ifndef PASTA_HPP
#define PASTA_HPP

inline int pasta2seq(int inp) {
  switch (inp) {
    case 'a': return 'a';
    case 'c': return 'c';
    case 'g': return 'g';
    case 't': return 't';
    case 'n': return 'n';

    case 'A': return 'n';
    case 'C': return 'n';
    case 'G': return 'n';
    case 'T': return 'n';

    case '~': return 'c';
    case '?': return 'g';
    case '@': return 't';

    case '=': return 'a';
    case ':': return 'g';
    case ';': return 't';

    case '#': return 'a';
    case '&': return 'c';
    case '%': return 't';

    case '*': return 'a';
    case '+': return 'c';
    case '-': return 'g';

    case 'Q': return 'a';
    case 'S': return 'c';
    case 'W': return 'g';
    case 'd': return 't';
    case 'Z': return 'n';

    case '!': return 0;
    case '$': return 0;
    case '7': return 0;
    case 'E': return 0;
    case 'z': return 0;

    case '\'': return 'a';
    case '"': return 'c';
    case ',': return 'g';
    case '_': return 't';

    default: return -1;

  }

  return -1;
}

inline int pasta2ref(int inp) {

  switch (inp) {
    case 'a': return 'a';
    case 'c': return 'c';
    case 'g': return 'g';
    case 't': return 't';
    case 'n': return 'n';

    case 'A': return 'a';
    case 'C': return 'c';
    case 'G': return 'g';
    case 'T': return 't';

    case '~': return 'a';
    case '?': return 'a';
    case '@': return 'a';

    case '=': return 'c';
    case ':': return 'c';
    case ';': return 'c';

    case '#': return 'g';
    case '&': return 'g';
    case '%': return 'g';

    case '*': return 't';
    case '+': return 't';
    case '-': return 't';

    case 'Q': return 0;
    case 'S': return 0;
    case 'W': return 0;
    case 'd': return 0;
    case 'Z': return 0;

    case '!': return 'a';
    case '$': return 'c';
    case '7': return 'g';
    case 'E': return 't';
    case 'z': return 'n';

    case '\'': return 'n';
    case '"': return 'n';
    case ',': return 'n';
    case '_': return 'n';

    default: return -1;

  }

  return -1;

}

inline int pasta_convert(int ref, int inp) {

  if (ref=='a') {
    if      (inp=='a') { return 'a'; }
    else if (inp=='c') { return '~'; }
    else if (inp=='g') { return '?'; }
    else if (inp=='t') { return '@'; }
    else if (inp=='n') { return 'A'; }
    else if (inp==0)   { return '!'; }
  }

  else if (ref=='c') {
    if      (inp=='a') { return '='; }
    else if (inp=='c') { return 'c'; }
    else if (inp=='g') { return ':'; }
    else if (inp=='t') { return ';'; }
    else if (inp=='n') { return 'C'; }
    else if (inp==0)   { return '$'; }
  }

  else if (ref=='g') {
    if      (inp=='a') { return '#'; }
    else if (inp=='c') { return '&'; }
    else if (inp=='g') { return 'g'; }
    else if (inp=='t') { return '%'; }
    else if (inp=='n') { return 'G'; }
    else if (inp==0)   { return '7'; }
  }

  else if (ref=='t') {
    if      (inp=='a') { return '*'; }
    else if (inp=='c') { return '+'; }
    else if (inp=='g') { return '-'; }
    else if (inp=='t') { return 't'; }
    else if (inp=='n') { return 'T'; }
    else if (inp==0)   { return 'E'; }
  }

  else if (ref=='n') {
    if      (inp=='a') { return '\''; }
    else if (inp=='c') { return '"'; }
    else if (inp=='g') { return ','; }
    else if (inp=='t') { return '_'; }
    else if (inp=='n') { return 'n'; }
    else if (inp==0)   { return 'z'; }
  }

  else if (ref==0) {
    if      (inp=='a') { return 'Q'; }
    else if (inp=='c') { return 'S'; }
    else if (inp=='g') { return 'W'; }
    else if (inp=='t') { return 'd'; }
    else if (inp=='n') { return 'Z'; }
    else if (inp==0)   { return 0; }
  }

  return 0;
}



#endif
