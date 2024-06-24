#!/bin/sh

curl "http://api:$HTTP_API_PORT/download?id=$UPLOAD_ID" -J --output-dir /tmp/ -O