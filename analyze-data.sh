#!/usr/bin/env bash

# finds all json files from data dir and applyes the bellow filters, until new monkey will fix the need for this :) 
find ./data/ -type f -name "*.json" | while read file; do \
printf " >>> file: %s\n" $file; \
jq '."url"' $file; \
jq '."timestamp"' $file ; \
jq '."stats"' $file; \
jq '.["data"]' $file | grep status | sort | uniq -c; done
