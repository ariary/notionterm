package notionterm

import (
	"context"
	"fmt"
	"strings"

	stringslice "github.com/ariary/go-utils/pkg/stringSlice"
	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

const CONFIGURATION string = "Configuration"
const TERMINAL string = "Terminal"
const TARGET = "Target"
const PORT = "Port"

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

//RequestButtonBlock: retrieve "button" widget (embed block)
func RequestButtonBlock(client *notionapi.Client, pageid string) (terminal notionapi.EmbedBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return terminal, err
	}
	return GetButtonBlock(children)
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

//RequestTerminalBlock: retrieve "terminal" block (code blocks)
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
			URL: url,
		},
	}

	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Embed: &widget.Embed,
	}

	return client.Block.Update(context.Background(), buttonID, updateReq)
}

//UpdateButtonCaption: update caption of the given button widget
func UpdateButtonCaption(client *notionapi.Client, button notionapi.EmbedBlock, caption string) (notionapi.Block, error) {
	//construct code block containing request
	widget := button

	captionRich := notionapi.RichText{
		Type: notionapi.ObjectTypeText,
		Text: notionapi.Text{
			Content: caption,
		},
		Annotations: &notionapi.Annotations{
			Bold:   false,
			Italic: true,
			Code:   true,
			Color:  "green",
		},
	}

	widget.Embed.Caption = []notionapi.RichText{captionRich}
	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Embed: &widget.Embed,
	}

	return client.Block.Update(context.Background(), button.ID, updateReq)
}

//UpdateCodeCaption: update caption of the given code block
// func UpdateCodeCaption(client *notionapi.Client, code notionapi.CodeBlock, caption string) (notionapi.Block, error) {
// 	//construct code block containing request
// 	terminal := code

// 	captionRich := notionapi.RichText{
// 		Type: notionapi.ObjectTypeText,
// 		Text: notionapi.Text{
// 			Content: caption,
// 		},
// 		Annotations: &notionapi.Annotations{
// 			Bold:   false,
// 			Italic: true,
// 			Code:   true,
// 			Color:  "green",
// 		},
// 	}

// 	terminal.Code.RichText.Caption = []notionapi.RichText{captionRich}
// 	// send update request
// 	updateReq := &notionapi.BlockUpdateRequest{
// 		Code: &terminal.Code,
// 	}

// 	return client.Block.Update(context.Background(), code.ID, updateReq)
// }

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
	var nRich []notionapi.RichText
	if len(content) > 2000 { //Add multiple chunked richtext
		chunks := stringslice.ChunksString(content, 350)
		for i := 0; i < len(chunks); i++ {
			newLine := notionapi.RichText{
				Type: notionapi.ObjectTypeText,
				Text: notionapi.Text{
					Content: chunks[i],
				},
			}
			if i == 0 {
				nRich = append(rich, newLine)
			} else {
				nRich = append(nRich, newLine)
			}
		}
	} else {
		newLine := notionapi.RichText{
			Type: notionapi.ObjectTypeText,
			Text: notionapi.Text{
				Content: content,
			},
		}
		nRich = append(rich, newLine)
	}

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

	return AddRichText(client, codeBlock, "$")
}

//GetTableBlock: retrieve table block
func GetTableBlock(children notionapi.Blocks) (table notionapi.TableBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeTableBlock {
			table = *children[i].(*notionapi.TableBlock)
			return table, nil
		}
	}
	err = fmt.Errorf("failed retrieving table block")
	return table, err
}

//RequestTableBlock: retrieve table block by requetsing it
func RequestTableBlock(client *notionapi.Client, pageid string) (table notionapi.TableBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return table, err
	}
	return GetTableBlock(children)
}

//GetTableRowBlock: retrieve table row block
func GetTableRowBlock(children notionapi.Blocks) (tableRow notionapi.TableRowBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeTableRowBlock {
			if i > 0 && children[i-1].GetType() == notionapi.BlockTypeHeading3 {
				heading := *children[i-1].(*notionapi.Heading3Block)
				if strings.Contains(heading.Heading3.RichText[0].Text.Content, CONFIGURATION) {
					tableRow = *children[i].(*notionapi.TableRowBlock)
					return tableRow, nil
				}
			}
		}
	}
	err = fmt.Errorf("failed retrieving table row block")
	return tableRow, err
}

//RequestTableRowBlock: retrieve table row block by requetsing it
func RequestTableRowBlock(client *notionapi.Client, pageid string) (tableRow notionapi.TableRowBlock, err error) {

	tableBlock, err := RequestTableBlock(client, pageid)
	if err != nil {
		return tableRow, err
	}

	tableBlockChildren, err := client.Block.GetChildren(context.Background(), tableBlock.ID, nil)
	if err != nil {
		return tableRow, err
	}

	return GetTableRowBlock(tableBlockChildren.Results)
}

//GetTableRowBlockbyHeader: retrieve table row block providing its header value
func GetTableRowBlockbyHeader(children notionapi.Blocks, header string) (tableRow notionapi.TableRowBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeTableRowBlock {
			tableRowTmp := *children[i].(*notionapi.TableRowBlock)
			//check config is above

			if len(tableRowTmp.TableRow.Cells) < 0 {
				continue
			}
			if tableRowTmp.TableRow.Cells[0][0].Text.Content == header {
				return tableRowTmp, nil
			}

		}
	}
	err = fmt.Errorf("Failed retrieving table row block")
	return tableRow, err
}

//RequestTableRowBlock: retrieve table row block providing its header valueby requetsing it
func RequestTableRowBlockByHeader(client *notionapi.Client, pageid string, header string) (tableRow notionapi.TableRowBlock, err error) {

	tableBlock, err := RequestTableBlock(client, pageid)
	if err != nil {
		return tableRow, err
	}

	tableBlockChildren, err := client.Block.GetChildren(context.Background(), tableBlock.ID, nil)
	if err != nil {
		return tableRow, err
	}

	return GetTableRowBlockbyHeader(tableBlockChildren.Results, header)
}

func RequestRowValueByHeader(client *notionapi.Client, pageid string, header string) (result string, err error) {
	tableRow, err := RequestTableRowBlockByHeader(client, pageid, header)
	if err != nil {
		return "", err
	}
	if len(tableRow.TableRow.Cells) < 2 {
		err = fmt.Errorf("failed retrieving value in table row (seems that the row does not have more than 1 columns)")
		return "", err
	} else if len(tableRow.TableRow.Cells[1]) < 1 {
		err = fmt.Errorf("failed retrieving value in table row (seems that the value is empty)")
		return "", err
	}
	result = tableRow.TableRow.Cells[1][0].Text.Content
	return result, err
}

//RequestTargetUrlFromConfig: return the value of the cell specifying the target urll/ip
func RequestTargetUrlFromConfig(client *notionapi.Client, pageid string) (targetUrl string, err error) {
	return RequestRowValueByHeader(client, pageid, TARGET)
}

//RequestTargetUrl: return the value of the cell specifying the target urll/ip
func RequestPortFromConfig(client *notionapi.Client, pageid string) (port string, err error) {

	return RequestRowValueByHeader(client, pageid, PORT)
}
