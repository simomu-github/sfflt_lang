#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname $0) && pwd)

declare -A EXAMPLES=(
	# test_file_name expect_output
	['hello_world']='Hello World!'
	['fib']='55'
	['local_var_and_logiral_operation']='11'
	['recursion_local_variable']='321'
	['stable_sort']='[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]'
)

has_failure=false

for example in "${!EXAMPLES[@]}"; do
	file=$example
	expect=${EXAMPLES[$example]}

	echo "Run ${file}"
	go run $SCRIPT_DIR/../cmd/sfflt_lang.go -format pretty -output $SCRIPT_DIR/test.fflt $SCRIPT_DIR/$file.sflt
	actual=$($FFLT_LANG $SCRIPT_DIR/test.fflt)

	if [ "$expect" == "$actual" ]; then
		echo "File: ${file} OK"
	else
		echo "File: ${file} expect: ${expect}, actual: ${actual}"
		has_failure=true
	fi
done

if [ "$has_failure" = true ]; then
	exit 1
fi
