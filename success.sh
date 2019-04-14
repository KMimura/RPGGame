#! /bin/bash

curl -X POST -H "Content-Type: application/json" -d '{"status":0,"url":"'"$TRAVIS_JOB_WEB_URL"'","message":"'"$TRAVIS_COMMIT_MESSAGE"'"}' $API_GATEWAY