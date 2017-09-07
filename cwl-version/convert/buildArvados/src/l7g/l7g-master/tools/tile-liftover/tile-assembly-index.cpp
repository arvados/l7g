// g++ tile-assembly-index.cpp -o tile-assembly-index
//
#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>
#include <errno.h>

#include <string>
#include <vector>

int main(int argc, char **argv) {
  int i, j, k, n, ch;
  std::string fn="-";
  FILE *fp=stdin;

  std::vector< std::string > tok_v;
  std::string line, name;
  int read_state=0;
  int cur_line_count=0, cur_char_count=0, cur_line_len=0;
  int cur_char_offset=0;
  int tot_line_count=0;

  if (argc<2) {
    printf("provide assembly file to index\n");
    exit(1);
  }

  fn = argv[1];

  if (fn!="-") {
    fp = fopen(fn.c_str(), "r");
    if (!fp) { perror(fn.c_str()); exit(1); }
  }

  while (!feof(fp)) {
    ch = fgetc(fp);
    if (ch==EOF) { continue; }
    if (ch=='\n') {
      tot_line_count++;

      if (line.size()==0) {
        fprintf(stderr, "zero line lengths not allow on char %i, line %i\n", cur_char_offset, tot_line_count);
        exit(1);
      }

      if (line[0]=='>') {

        if (name.size()>0) {
          printf("%s\t%i\t%i\t%i\t%i\n",
              name.c_str(),
              cur_char_count,
              cur_char_offset,
              cur_line_len,
              cur_line_len+1);
        }

        cur_char_offset += cur_line_count*(cur_line_len+1);

        name.clear();
        for (i=1; i<line.size(); i++) { name += line[i]; }
        read_state = 1;
        cur_line_count=0;
        cur_char_count=0;
        cur_line_len=0;

        cur_char_offset += (int)line.size()+1;
      }
      else {
        cur_line_count++;
        if (cur_line_len==0) { cur_line_len = (int)line.size(); }
        if (cur_line_len != (int)line.size()) {
          fprintf(stderr, "line mismatch on %i (char offset %i), expecting %i, got %i\n",
              tot_line_count,
              cur_char_offset,
              cur_line_len,
              (int)line.size());
          exit(3);
        }
        cur_char_count += cur_line_len+1;
      }

      line.clear();

      continue;
    }

    line += (char)ch;

  }

  if (name.size()>0) {
    printf("%s\t%i\t%i\t%i\t%i\n",
        name.c_str(),
        cur_char_count,
        cur_char_offset,
        cur_line_len,
        cur_line_len+1);
  }


  if (fp!=stdin) { fclose(fp); }
}
