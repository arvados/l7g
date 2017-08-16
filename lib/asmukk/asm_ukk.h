#ifndef ASM_UKK_H
#define ASM_UKK_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define ASM_UKK_MISMATCH 3
#define ASM_UKK_GAP 2

int asm_ukk_score(char *, char *);
int asm_ukk_score2(char *, char *, int, int);
int asm_ukk_score3(char *a, char *b, int (*score_func)(char,char));

int asm_ukk_align(char **, char **, char *, char *);
int asm_ukk_align2(char **X, char **Y, char *a, char *b, int mismatch, int gap, char gap_char);
int asm_ukk_align3(char **X, char **Y, char *a, char *b, int (*score_func)(char, char), char gap_char);


int sa_align_ukk(char **, char **, char *, char *, int);
int sa_align_ukk2(char **, char **, char *, char *, int, int, int, char);
int sa_align_ukk3(char **X, char **Y, char *a, char *b, int threshold, int (*score_func)(char, char), char gap_char);

#endif
