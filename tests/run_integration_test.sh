#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname $0) && pwd)
EXAMPLES=(
	# test_file_name expect_output
	"hello_world,Hello World!"
	"fib,55"
	"local_var_and_logiral_operation,11"
	"recursion_local_variable,321"
)

for example in "${EXAMPLES[@]}"; do
	file=$(echo "$example" | cut -d',' -f1)
	expect=$(echo "$example" | cut -d',' -f2)

	echo "Run ${file}"
	go run $SCRIPT_DIR/../cmd/sfflt_lang.go -format pretty -output $SCRIPT_DIR/test.fflt $SCRIPT_DIR/$file.sflt
	actual=$($FFLT_LANG $SCRIPT_DIR/test.fflt)

	if [ "$expect" == "$actual" ]; then
		echo "File: ${file} OK"
	else
		echo "File: ${file} expect: ${expect}, actual: ${actual}"
	fi
done
