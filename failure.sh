#! /bin/bash

curl -X POST -H "Content-Type: application/json" -d '{"status":1,"id":"'"$TRAVIS_COMMIT"'"}' $API_GATEWAY
