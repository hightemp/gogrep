#!/bin/bash

echo "Test 1"
echo "gogrep"
time ./gogrep "Placerat" test/file1.txt
echo "grep"
time grep "Placerat" test/file1.txt

echo "Test 2"
echo "gogrep"
time ./gogrep "Richard" test/file1.txt test/file2.txt
echo "grep"
time grep "Richard" test/file1.txt test/file2.txt

echo "Test 3"
echo "gogrep"
time ./gogrep -r "Richard" ./test
echo "grep"
time grep -r "Richard" ./test