#!/bin/bash

if [[ $# -ne 3 ]]; then
    echo "usage: ./static-build.sh \$NOTIONTERM_PAGE_URL \$NOTION_TOKEN \$GOOS"
    exit 92
fi


export URL=$1
export TOKEN=$2

GOOS=$3 go build  -ldflags "-X 'main.PageUrl=$URL' -X 'main.Token=$TOKEN'" notionterm.go