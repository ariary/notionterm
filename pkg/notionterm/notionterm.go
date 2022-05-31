package notionterm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

func CreateHeadingCaptionBlock(config Config) notionapi.BlockID {
	resp, err := AppendHeadingBlock(config.Client, config.PageID, "📁")
	if err != nil {
		fmt.Println("❌ Failed to create heading caption block:", err)
		os.Exit(92)
	} else if len(resp.Results) < 1 {
		fmt.Println("❌ Failed to retrieve heading id after creation:", err)
		os.Exit(92)
	}

	return resp.Results[0].GetID()
}

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

// CreateTerminalBlock: create the block code simulating a terminal
func CreateTerminalBlock(config Config) notionapi.BlockID {
	resp, err := AppendCodeBlock(config.Client, config.PageID, config.PS1)
	if err != nil {
		fmt.Println("❌ Failed to create code block:", err)
		os.Exit(92)
	} else if len(resp.Results) < 1 {
		fmt.Println("❌ Failed to retrieve terminal code bock id after creation:", err)
		os.Exit(92)
	}
	return resp.Results[0].GetID()
}

//DeleteEmbed: find the last embed/bookmark block of a page and delete it
func DeleteEmbed(config Config) (err error) {
	children, err := notionion.RequestProxyPageChildren(config.Client, config.PageID)
	if err != nil {
		return err
	}
	//Check if last embed is well loaded
	embed := children[len(children)-1]
	if embed.GetType() == notionapi.BlockTypeBookmark || embed.GetType() == notionapi.BlockTypeEmbed {
		//delete bookmark
		if _, err := config.Client.Block.Delete(context.Background(), embed.GetID()); err != nil {
			fmt.Println("Failed deleting embed/bookmark:", err)
			os.Exit(92)
		}
	} else {
		err = errors.New("last block is not an embed or bookmark block")
	}

	return err
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
			if strings.Contains(cmd, "\n") && strings.HasPrefix(cmd, config.PS1) {
				if isCommand(cmd) {
					cmdSplit := strings.Split(cmd, config.PS1)
					if len(cmdSplit) > 1 {
						cmd = cmdSplit[1]
					}
					cmd = strings.Replace(cmd, "\n", "", -1)
					if !handleSpecialCommand(config, termBlock, cmd) {
						//Execute it and print
						ExecAndPrint(config, termBlock, cmd)
					}

					//refresh+add new terminal line ($)
					termBlock, err = RequestTerminalBlock(config.Client, config.PageID, config.TerminalBlockId)
					if err != nil {
						fmt.Println(err)
						continue
					}
					AddTermLine(config, termBlock)
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

//handleSpecialCommand: make actions for specific  command
func handleSpecialCommand(config *Config, termBlock notionapi.CodeBlock, cmd string) (isSpecial bool) {
	if strings.HasPrefix(cmd, "cd ") { //TODO handle cd without space ("cd " = "cd")
		//change path
		cmdSplit := strings.Split(cmd, " ")
		if len(cmdSplit) > 1 {
			path := cmdSplit[1]
			//check path
			if path == "" {
				if user, err := user.Current(); err == nil {
					path = user.HomeDir
				}
			} else if !strings.HasPrefix(path, "/") {
				path = config.Path + "/" + path
				if pathTmp, err := filepath.Abs(path); err == nil {
					path = pathTmp
				}
			}
			if info, err := os.Stat(path); !os.IsNotExist(err) && info.IsDir() {
				//update caption
				if _, err := UpdateCaptionById(config.Client, config.PageID, config.CaptionBlock, "📁 "+path); err != nil {
					fmt.Println("Failed updating caption with current path:", err)
				} else {
					config.Path = path
					fmt.Println("📁 Change directory:", path)
				}
			}
		} else {
			fmt.Println("Failed retrieving directory in 'cd' command:", cmd)
		}
		return true
	} else if strings.HasPrefix(cmd, "clear") {
		fmt.Println("🦆 Clear terminal")
		UpdateCodeContent(config.Client, termBlock.ID, "")
		return true
	} else if strings.HasPrefix(cmd, "bye") {
		fmt.Println("👋 Close notionterm")
		//UpdateCodeContent(config.Client, termBlock.ID, "👋 see u")
		if _, err := AddRichText(config.Client, termBlock, "👋 see u"); err != nil {
			fmt.Println("failed to add rich text in terminal code block:", err)
		}
		os.Exit(0)
		return true
	}
	return false
}

//ExecAndPrint: execute command and print the result in code block
func ExecAndPrint(config *Config, termBlock notionapi.CodeBlock, cmd string) {
	fmt.Println("📟", cmd)
	var flag string
	switch config.Shell {
	case "cmd.exe":
		flag = "/C"
	default:
		flag = "-c"
	}
	commandExec := exec.Command(config.Shell, flag, cmd)
	commandExec.Dir = config.Path
	stdout, err := commandExec.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Print the output
	//fmt.Println(string(stdout))
	if _, err := AddRichText(config.Client, termBlock, string(stdout)); err != nil {
		fmt.Println("failed to add rich text in terminal code block:", err)
	}
}
