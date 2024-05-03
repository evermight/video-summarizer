#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/../.env
source $SCRIPT_DIR/save.sh

curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $GPT_API_KEY" \
  -d @$SCRIPT_DIR/conversations/$topic > $SCRIPT_DIR/response.json

gptresponse=$(cat $SCRIPT_DIR/response.json | jq -r '.choices | map(select(.message.role == "assistant")) | .[0].message.content');
echo $gptresponse;

$SCRIPT_DIR/save.sh -t "$topic" -c "$gptresponse" -r assistant
rm -rf $SCRIPT_DIR/response.json
