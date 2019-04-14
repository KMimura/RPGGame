#! /bin/bash

curl -X POST -H "Content-Type: application/json" -d '{"status":0,"id":"'"$TRAVIS_COMMIT"'","message":"'"$TRAVIS_COMMIT_MESSAGE"'"}' $API_GATEWAY