package notionterm

import (
	"fmt"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

func NotionTerm(client *notionapi.Client, pageid string, play chan struct{}, pause chan struct{}) {
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-pause:
			//fmt.Println("pause")
			select {
			case <-play:
				//fmt.Println("play")
			}
		default:
			cmd, err := RequestTerminalCodeContent(client, pageid)
			if err != nil {
				fmt.Println(err)
			}
			if strings.Contains(cmd, "\n") {
				if isCommand(cmd) {
					cmd = strings.Split(cmd, "$ ")[1] //todo check len
					//execute it
					fmt.Println(cmd)
				}
			}
		}
	}
}

//check if a command really is
func isCommand(command string) bool {
	if command[len(command)-2] == '\\' {
		return false
	}
	return true
}
