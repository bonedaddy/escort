#! /bin/bash

DELFATE_OUT=""
BASE64_OUT=""
BUFFER=""
for i in $(seq 1 10000); do
    BUFFER+="$i"
    # compress and base64 encode buffer
    DEFLATE_OUT=$(./escort --input "$BUFFER" compress)
    # just base64 encode buffer
    BASE64_OUT=$(./escort --input "$BUFFER" base64 encode)
    # get length of deflate + base64
    DEFLATE_OUT_LEN=$(echo "$DEFLATE_OUT" | wc -c)
    # get length of base64
    BASE64_OUT_LEN=$(echo "$BASE64_OUT" | wc -c)
    if [[ "$DEFLATE_OUT_LEN" -lt "$BASE64_OUT_LEN" ]]; then
        echo "minimum input characters for deflate + base64 to be smaller than just base64 is $DEFLATE_OUT_LEN"
        exit 0
    fi
done