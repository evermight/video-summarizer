#!/bin/bash
while getopts i: option
do
    case "${option}" in
        i) id=${OPTARG};;
    esac
done

if [ -z $id ]; then
  echo '-i for video id required'
  exit
fi

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source $SCRIPT_DIR/../.env
#curl -X GET -H "Authorization: Bearer $GOOGLE_ACCESS_TOKEN" "https://www.googleapis.com/youtube/v3/captions/?key=$GOOGLE_API_KEY&part=id&videoId=$id"
curl -X GET -H "Authorization: Bearer $GOOGLE_ACCESS_TOKEN" "https://www.googleapis.com/youtube/v3/captions/?key=$GOOGLE_API_KEY&part=snippet&videoId=$id"
