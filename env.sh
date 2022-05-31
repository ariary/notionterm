#!/bin/bash

echo -n "Enter integration token: " 
read -s token
echo
echo -n "Enter notion page url: "
read pageurl

export NOTION_TOKEN="$token"
export NOTION_PAGE_URL="$pageurl"