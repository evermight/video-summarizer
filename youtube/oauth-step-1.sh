#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/../.env

scopes=("profile" "https://www.googleapis.com/auth/youtubepartner" "https://www.googleapis.com/auth/youtube.force-ssl")
scopelist=""
for value in "${scopes[@]}"
do
  scopelist="$scopelist$value%20"
done

echo "Visit this url in a web browser:"
echo "";
echo "https://accounts.google.com/o/oauth2/auth?client_id=$GOOGLE_CLIENT_ID&redirect_uri=http://localhost&scope=$scopelist&email&response_type=code&include_granted_scopes=true&access_type=offline&state=state_parameter_passthrough_value";
echo "";
echo "Once you have been redirected to http://localhost, copy the <CODE> from &code=<CODE> in the url";
echo "";
