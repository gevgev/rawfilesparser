#!/bin/bash

if [ "$#" -lt 1 ]; then
  echo "ERROR: missing required parameters:  <instances> "
  echo "|"
  echo "| instances:   1, 2, 0 (= destroy)"
  exit 1
fi

readonly instances=$1

if [ "$instances" -gt 2 -o "$instances" -lt 0 ]; then
  echo "ERROR: Allowed instances count is 0, 1, 2"
  exit 1
fi

echo "Instances requested: $instances"

terraform plan -var "instances=$instances"

plan_result=$?

if [ "$plan_result" -ne 0 ] ; then
	echo "Errors encountered. Stopping"
	exit "$plan_result"
fi

echo "Is this what you want to do? (only 'yes' will be accepted as aknowledgement)"

read confirm


if [ "$confirm" = "yes" ]; then
	echo "Proceeding with the requested update"
	terraform apply -var "instances=$instances"
else
	echo "Leaving as is"
fi
