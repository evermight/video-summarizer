#!/bin/bash

while getopts t:c:r: option
do
    case "${option}" in
        t) topic=${OPTARG};;
        c) content=${OPTARG};;
        r) role=${OPTARG};;
    esac
done

if [ -z "$topic" ]; then
	echo '-t for topic required (letters and numbers only)'
  exit
fi
if [ -z "$content" ]; then
	echo '-c for content required'
  exit
fi
if [ -z "$role" ]; then
        role="user"
fi

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

isvalid=0
validroles=("user" "system" "assistant")
for value in "${validroles[@]}"
do
  scopelist="$scopelist$value%20"
  if [ "$value" == "$role" ]; then
    isvalid=1
  fi
done

if [ "$isvalid" == "0" ]; then
  echo "Invalid -r for role: can only be user, system or assistant"
  exit;
fi

convdir="$SCRIPT_DIR/conversations"
mkdir -p $convdir;

convfile="$convdir/$topic"

if [ ! -f "$convdir/$topic" ]; then
    echo '{"model": "gpt-3.5-turbo-0125", "messages": []}' > "$convdir/$topic"
fi
message=$(jq --null-input \
  --arg role "$role" \
  --arg content "$content" \
  '{"role": $role, "content": $content}')

convcontent=$(cat $convfile | jq '.messages += ['"$message"']');
echo $convcontent > $convfile;
