#!/bin/bash

set -x

if [ "$#" -lt 1 ]; then
  echo "ERROR: missing required parameters:  <instances> "
  echo "|"
  echo "| instances:   1, 2, 0 (= destroy)"
  exit 1
fi

readonly instances=$1
readonly s3_bucket="r31analysis/click" #$2
readonly ftp_folder="/iGuide/Click" #$3
readonly ftp_file="Click-Tacoma-test.zip" #$4

if [ "$instances" -gt 2 -o "$instances" -lt 0 ]; then
  echo "ERROR: Allowed instances count is 0, 1, 2"
  exit 1
fi

echo "Instances requested: $instances"

terraform plan -var "instances=$instances" -var "s3-bucket='$s3_bucket'" -var "ftp-folder='$ftp_folder'" -var "ftp-file='$ftp_file'"

plan_result=$?

if [ "$plan_result" -ne 0 ] ; then
	echo "Errors encountered. Stopping"
	exit "$plan_result"
fi

echo "Is this what you want to do? (only 'yes' will be accepted as aknowledgement)"

read confirm


if [ "$confirm" = "yes" ]; then
	echo "Proceeding with the requested update"
	terraform apply -var "instances=$instances" -var "s3-bucket='$s3_bucket'" -var "ftp-folder='$ftp_folder'" -var "ftp-file='$ftp_file'"
else
	echo "Leaving as is"
fi
