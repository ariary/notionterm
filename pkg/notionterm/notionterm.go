package notionterm

import (
	"fmt"
	"os/exec"
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
			termBlock, err := RequestTerminalBlock(client, pageid)
			if err != nil {
				fmt.Println(err)
				continue
			}
			cmd, err := GetTerminalLastRichText(termBlock)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("last:", cmd)
			if strings.Contains(cmd, "\n") && strings.HasPrefix(cmd, "$ ") {
				if isCommand(cmd) {
					cmdSplit := strings.Split(cmd, "$ ")
					if len(cmdSplit) > 1 {
						cmd = cmdSplit[1] //todo check len
					}
					//execute it and print
					fmt.Println(cmd)
					cmmandExec := exec.Command("sh", "-c", cmd)
					stdout, err := cmmandExec.Output()

					if err != nil {
						fmt.Println(err.Error())
						return
					}
					// Print the output
					//fmt.Println(string(stdout))
					if _, err := AddRichText(client, termBlock, string(stdout)); err != nil {
						fmt.Println(err)
					}

					//refresh+add new terminal line ($)
					termBlock, err = RequestTerminalBlock(client, pageid)
					if err != nil {
						fmt.Println(err)
						continue
					}
					AddTermLine(client, termBlock)
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
