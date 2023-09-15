#!/bin/sh

# script from https://github.com/rui314/chibicc/blob/master/test.sh

assert() {
  expected="$1"
  input="$2"

  go run gosling.go "$input" > test.s || exit
  clang -o test test.s
  ./test
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assert 0 0
assert 42 42

echo OK