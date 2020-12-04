#!/usr/bin/env bash
set -xu

go build -o go-cp

result=$(./go-cp -from testdata/input.txt -to ./testdata)
if [ "${result}" != 'Error: directories are not supported' ]; then
        echo "unexpected output"
        exit 1;
fi

result=$(./go-cp -from /dev/urandom -to out.txt)
if [ "${result}" != 'Error: unsupported file' ]; then
        echo "unexpected output"
        exit 1;
fi

result=$(./go-cp -from testdata/input.txt -to testdata/input.txt)
if [ "${result}" != 'Error: can not copy file into itself' ]; then
        echo "unexpected output"
        exit 1;
fi

result=$(./go-cp -from input.txt -to out.txt)
if [ "${result#*txt:}" != ' no such file or directory' ]; then
        echo "unexpected output"
        exit 1;
fi


set -xeuo pipefail

./go-cp -from testdata/input.txt -to out.txt
cmp out.txt testdata/out_offset0_limit0.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10
cmp out.txt testdata/out_offset0_limit10.txt

./go-cp -from testdata/input.txt -to out.txt -limit 1000
cmp out.txt testdata/out_offset0_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10000
cmp out.txt testdata/out_offset0_limit10000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000
cmp out.txt testdata/out_offset100_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
cmp out.txt testdata/out_offset6000_limit1000.txt

rm -f go-cp out.txt
echo "PASS"
