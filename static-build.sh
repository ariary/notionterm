#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "usage: ./static-build.sh \$NOTIONTERM_PAGE_URL \$NOTION_TOKEN"
    exit 92
fi


export URL=$1
export TOKEN=$2

go build  -ldflags "-X 'github.com/ariary/notionterm/pkg/notionterm.pageurl=$URL' -X 'github.com/ariary/notionterm/pkg/notionterm.token=$TOKEN'" notionterm.go