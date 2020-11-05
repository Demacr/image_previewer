#!/bin/bash

curl localhost:8080/fill/200/300/nginx/image.jpg -o image.jpg

rm -rf image.jpg
