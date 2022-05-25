package notionterm

import (
	"fmt"
	"time"
)

var pageurl, token string

//V2

func CreateButtonBlock() {
	fmt.Println("Create button")
}

func CreateTerminalBlock() {
	fmt.Println("Create terminal")
}

func DeleteEmbed() {
	fmt.Println("delete embed")
}

//NotionTerm: "Infinite loop" to read the content of terminal code block and execute it if it is a command, then returning stdout
func NotionTermV2(config Config, play chan struct{}, pause chan struct{}) {
	fmt.Println("Notiontermv2")
	for {
		time.Sleep(config.Delay)
		select {
		case <-pause:
			fmt.Println("pause")
			select {
			case <-play:
				fmt.Println("play")
			}
		default:
			fmt.Println("default")

		}
	}
}
