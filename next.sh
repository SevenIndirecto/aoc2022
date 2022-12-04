#!/usr/bin/env bash

set -e

# import .env
set -o allexport; source .env; set +o allexport
if [[ -z $YEAR || -z $SESSION ]]; then
  echo "Set YEAR and SESSION in an .env file"
  exit
fi

if [[ ! -d 'template' ]]; then
  echo "Make sure you're in the root directory and template dir exists".
  exit
fi

# get previous day dir name
PREV_DAY=$(ls -d */ | grep -E "[0-9]{2}" | sed 's#/##' | tail -n1)

if [[ -z $PREV_DAY ]]; then
  NEXT_NUM=1
else
  # 10# to force base 10, since our format can be "08" and interpreted as octal
  NEXT_NUM=$((10#${PREV_DAY}+1))
fi

NEXT=$(printf "%02d" "$NEXT_NUM")
mkdir "$NEXT"
cp template/main.go "$NEXT/aoc.go"
cp template/main_test.go "$NEXT/aoc_test.go"
touch "$NEXT/test.txt"

curl "https://adventofcode.com/$YEAR/day/$NEXT_NUM/input" -s -H "cookie: session=$SESSION" > "$NEXT/input.txt"
echo "Day $NEXT_NUM ready, instructions at https://adventofcode.com/$YEAR/day/$NEXT_NUM"

