package notionterm

import (
	"context"
	"fmt"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

//GetButtonBlock: retrieve "button" block (embed blocks)
func GetButtonBlock(children notionapi.Blocks) (button notionapi.EmbedBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeEmbed {
			button = *children[i].(*notionapi.EmbedBlock)
			return button, nil
		}
	}
	err = fmt.Errorf("Failed retrieving \"button\" widget")
	return button, err
}

//GetTerminalBlock: retrieve "terminal" block (code blocks)
func GetTerminalBlock(children notionapi.Blocks) (terminal notionapi.CodeBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeCode {
			terminal = *children[i].(*notionapi.CodeBlock)
			//to do check if terminal is under the button
			return terminal, nil
		}
	}
	err = fmt.Errorf("Failed retrieving \"terminal\" section")
	return terminal, err
}

//GetTerminalBlock: retrieve "terminal" block (code blocks)
func RequestTerminalBlock(client *notionapi.Client, pageid string) (terminal notionapi.CodeBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return terminal, err
	}
	return GetTerminalBlock(children)
}

//RequestTerminalCodeContent: Obtain the content of code block object under the request heading
func RequestTerminalCodeContent(client *notionapi.Client, pageid string) (terminal string, err error) {

	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return "", err
	}
	return GetTerminalCodeContent(children)
}

//GeTerminalCodeContent: Obtain the content of code block object under the request heading whithout making request
func GetTerminalCodeContent(children notionapi.Blocks) (terminal string, err error) {
	termCode, err := GetTerminalBlock(children)
	if err != nil {
		return "", err
	}
	terminal = termCode.Code.RichText[0].PlainText
	return terminal, err
}

//GetTerminalLastRichText: Obtain the last RichText
func GetTerminalLastRichText(termCode notionapi.CodeBlock) (terminal string, err error) {
	terminal = termCode.Code.RichText[len(termCode.Code.RichText)-1].PlainText
	return terminal, err
}

//UpdateButtonUrl: update url of the button widget
func UpdateButtonUrl(client *notionapi.Client, buttonID notionapi.BlockID, url string) (notionapi.Block, error) {
	//construct code block containing request
	widget := notionapi.EmbedBlock{
		Embed: notionapi.Embed{
			Caption: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{
						Content: "",
					},
					Annotations: &notionapi.Annotations{
						Bold:          false,
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
						Code:          false,
						Color:         "",
					},
				},
			},
			URL: url,
		},
	}

	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Embed: &widget.Embed,
	}

	return client.Block.Update(context.Background(), buttonID, updateReq)
}

//UpdateCodeContent: update code block with content
func UpdateCodeContent(client *notionapi.Client, codeBlockID notionapi.BlockID, content string) (notionapi.Block, error) {
	//construct code block containing request
	code := notionapi.CodeBlock{
		Code: notionapi.Code{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{
						Content: content,
					},
					Annotations: &notionapi.Annotations{
						Bold:          false,
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
						Code:          false,
						Color:         "",
					},
				},
			},
			Language: "shell",
		},
	}

	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Code: &code.Code,
	}

	return client.Block.Update(context.Background(), codeBlockID, updateReq)
}

//AddRichText: Add rich text in code
func AddRichText(client *notionapi.Client, codeBlock notionapi.CodeBlock, content string) (notionapi.Block, error) {
	rich := codeBlock.Code.RichText
	newLine := notionapi.RichText{
		Type: notionapi.ObjectTypeText,
		Text: notionapi.Text{
			Content: content,
		},
	}
	nRich := append(rich, newLine)
	//construct code block containing request
	code := notionapi.CodeBlock{
		Code: notionapi.Code{
			RichText: nRich,
			Language: "shell",
		},
	}
	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Code: &code.Code,
	}

	return client.Block.Update(context.Background(), codeBlock.ID, updateReq)
}

//AddTermLine: Add rich text with a new line and "$"
func AddTermLine(client *notionapi.Client, codeBlock notionapi.CodeBlock) (notionapi.Block, error) {

	return AddRichText(client, codeBlock, "\n$")
}
