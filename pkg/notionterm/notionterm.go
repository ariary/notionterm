package notionterm

import (
	"fmt"
	"os"
	"time"

	"github.com/jomei/notionapi"
)

//V2

//CreateButtonBlock: create embed block with url button
func CreateButtonBlock(config Config) notionapi.BlockID {
	buttonUrl := config.ExternalUrl + "/button"
	resp, err := AppendEmbedBlock(config.Client, config.PageID, buttonUrl)
	if err != nil {
		fmt.Println("❌ Failed to create button widget:", err)
		os.Exit(92)
	} else if len(resp.Results) < 1 {
		fmt.Println("❌ Failed to retrieve button widget id after creation:", err)
		os.Exit(92)
	}
	return resp.Results[0].GetID()
}

func CreateTerminalBlock(config Config) notionapi.BlockID {
	fmt.Println("Create terminal")
	resp, err := AppendCodeBlock(config.Client, config.PageID, config.PS1)
	if err != nil {
		fmt.Println("❌ Failed to create button widget:", err)
		os.Exit(92)
	} else if len(resp.Results) < 1 {
		fmt.Println("❌ Failed to retrieve button widget id after creation:", err)
		os.Exit(92)
	}
	return resp.Results[0].GetID()
}

func DeleteEmbed() {
	fmt.Println("delete embed")
}

//NotionTerm: "Infinite loop" to read the content of terminal code block and execute it if it is a command, then returning stdout
func NotiontermRun(config *Config, play chan struct{}, pause chan struct{}) {
	for {
		time.Sleep(config.Delay)
		select {
		case <-pause:
			//fmt.Println("pause")
			select {
			case <-play:
				//fmt.Println("play")
			}
		default:
			// request last command
			termBlock, err := RequestTerminalBlock(config.Client, config.PageID, config.TerminalBlockId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			cmd, err := GetTerminalLastRichText(termBlock)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("last:", cmd)

		}
	}
}
