#!/usr/bin/env bash

packages=$(go list ./... \
  | grep -v "^github.com/twk/skeleton-go-api/cmd" \
  | grep -v "/mocks" \
)

go test $packages -covermode=atomic -coverprofile=coverage.out

EXPECTED_COVERAGE=${EXPECTED_COVERAGE:-80}
function die() {
  echo $*
  exit 1
}

cov=`go tool cover -func=coverage.out | tail -n 1 | awk '{print $3}' | sed 's/[^0-9\.]*//g'`
comparison=$(awk -v cov=$cov -v expected=$EXPECTED_COVERAGE 'BEGIN { print (cov >= expected) ? 1 : 0 }')

if [ "$comparison" -eq 1 ]; then
  echo "SUCCESS: Coverage is ~$cov% (minimum expected is $EXPECTED_COVERAGE%)"
else
  die "ERROR: Test coverage is not enough! Want at least $EXPECTED_COVERAGE% but only $cov% of tested packages are covered with tests."
fi
