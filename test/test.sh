#!/bin/bash

set -e

echo -n "-------> Simple case status code "
[ $(curl -I -L -s -w "%{http_code}" localhost:8080/fill/200/300/nginx/image.jpg -o image.jpg) -eq 200 ] && echo "succeded"
curl localhost:8080/fill/200/300/nginx/image.jpg -o image.jpg 2>/dev/null
echo -n "-------> Simple case image size "
[ $(file image.jpg | awk '{print $8;}') = "200x300," ] && echo "succeded" || echo "failed"
rm -rf image.jpg

echo -n "-------> Wrong name case status code "
[ $(curl -I -L -s -w "%{http_code}" localhost:8080/fill/200/300/nginx/image.jp -o image.jpg) -eq 500 ] && echo "succeded"
