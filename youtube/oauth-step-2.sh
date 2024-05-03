#!/bin/bash
while getopts c: option
do
    case "${option}" in
        c) code=${OPTARG};;
    esac
done

if [ -z $code ]; then
  echo '-c for <CODE> required'
  exit
fi

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/../.env

curl -X POST https://oauth2.googleapis.com/token -d "client_id=$GOOGLE_CLIENT_ID&client_secret=$GOOGLE_CLIENT_SECRET&redirect_uri=http://localhost&grant_type=authorization_code&code=$code" > $SCRIPT_DIR/tokens.json
echo "Credentials saved to $SCRIPT_DIR/tokens.json";
