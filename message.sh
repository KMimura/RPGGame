#! /bin/bash

curl -X POST -H "Content-Type: application/json" -d '{"text":"'"$TRAVIS_COMMIT"'"}' $API_GATEWAY