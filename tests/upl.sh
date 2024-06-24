#!/bin/sh

curl -v -i -X POST -H "Content-Type: multipart/form-data" -F "data=@/tmp/$TEST_UPL_FILE" http://api:$HTTP_API_PORT/upload