package notionterm

import (
	"context"
	"fmt"

	stringslice "github.com/ariary/go-utils/pkg/stringSlice"
	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

const CONFIGURATION string = "C O N F I G U R A T I O N"
const TERMINAL string = "T E R M I N A L"
const TARGET = "Target"
const PORT = "Port"
const SHELL = "Shell"

//GetButtonBlock: retrieve "button" block (embed blocks)
func GetButtonBlock(children notionapi.Blocks, buttonId notionapi.BlockID) (button notionapi.EmbedBlock, err error) {
	for i := len(children) - 1; i >= 0; i-- {
		if children[i].GetType() == notionapi.BlockTypeEmbed && children[i].GetID() == buttonId {
			button = *children[i].(*notionapi.EmbedBlock)
			return button, nil
		}
	}
	err = fmt.Errorf("Failed retrieving \"button\" widget")
	return button, err
}

//RequestButtonBlock: retrieve "button" widget (embed block)
func RequestButtonBlock(client *notionapi.Client, pageid string, buttonId notionapi.BlockID) (button notionapi.EmbedBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return button, err
	}
	return GetButtonBlock(children, buttonId)
}

//GetHeadingCaptionBlock: retrieve heading3 block by id.
func GetHeadingCaptionBlock(children notionapi.Blocks, headingID notionapi.BlockID) (headingCaption notionapi.Heading3Block, err error) {
	for i := len(children) - 1; i >= 0; i-- {
		if children[i].GetType() == notionapi.BlockTypeHeading3 && children[i].GetID() == headingID {
			headingCaption = *children[i].(*notionapi.Heading3Block)
			return headingCaption, nil
		}
	}
	err = fmt.Errorf("Failed retrieving caption heading")
	return headingCaption, err
}

//RequestButtonBlock: retrieve heading 3 block to store caption/path
func RequestHeadingCaptionBlock(client *notionapi.Client, pageid string, headingID notionapi.BlockID) (button notionapi.Heading3Block, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return button, err
	}
	return GetHeadingCaptionBlock(children, headingID)
}

//GetTerminalBlock: retrieve "terminal" block (code blocks)
func GetTerminalBlock(children notionapi.Blocks, termId notionapi.BlockID) (terminal notionapi.CodeBlock, err error) {
	for i := len(children) - 1; i >= 0; i-- {
		if children[i].GetType() == notionapi.BlockTypeCode && children[i].GetID() == termId {
			terminal = *children[i].(*notionapi.CodeBlock)
			//to do check if terminal is under the button
			return terminal, nil
		}
	}
	err = fmt.Errorf("Failed retrieving \"terminal\" section")
	return terminal, err
}

//RequestTerminalBlock: retrieve "terminal" block (code blocks)
func RequestTerminalBlock(client *notionapi.Client, pageId string, termId notionapi.BlockID) (terminal notionapi.CodeBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageId)
	if err != nil {
		return terminal, err
	}
	return GetTerminalBlock(children, termId)
}

//RequestTerminalCodeContent: Obtain the content of code block object under the request heading
func RequestTerminalCodeContent(client *notionapi.Client, pageid string, termId notionapi.BlockID) (terminal string, err error) {

	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return "", err
	}
	return GetTerminalCodeContent(children, termId)
}

//GeTerminalCodeContent: Obtain the content of code block object under the request heading whithout making request
func GetTerminalCodeContent(children notionapi.Blocks, termId notionapi.BlockID) (terminal string, err error) {
	termCode, err := GetTerminalBlock(children, termId)
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

//UpdateCaptionID: update caption of captionBlock. For the light mode the "caption block" is a heading3 so we change the content (not the caption)
func UpdateCaptionById(client *notionapi.Client, pageID string, captionBlock CaptionBlock, caption string) (notionapi.Block, error) {
	var updateReq *notionapi.BlockUpdateRequest
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

	switch captionBlock.Type {
	case notionapi.BlockTypeHeading3: //light mode
		if headingCaption, err := RequestHeadingCaptionBlock(client, pageID, captionBlock.Id); err != nil {
			return nil, err
		} else {
			headingCaption.Heading3.RichText = []notionapi.RichText{captionRich}
			// set update request
			updateReq = &notionapi.BlockUpdateRequest{
				Heading3: &headingCaption.Heading3,
			}
			client.Block.Update(context.Background(), captionBlock.Id, updateReq)
		}
	case notionapi.BlockTypeEmbed: //normal mode
		//construct code block containing request
		if button, err := RequestButtonBlock(client, pageID, captionBlock.Id); err != nil {
			return nil, err
		} else {
			button.Embed.Caption = []notionapi.RichText{captionRich}
			// set update request
			updateReq = &notionapi.BlockUpdateRequest{
				Embed: &button.Embed,
			}
		}
	}

	return client.Block.Update(context.Background(), captionBlock.Id, updateReq)
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
func AddTermLine(config *Config, codeBlock notionapi.CodeBlock) (notionapi.Block, error) {

	return AddRichText(config.Client, codeBlock, config.PS1)
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

//GetRowValue:return the value of the corresponding header in a table (does not check header only return second cell text content)
func GetRowValue(tableRow notionapi.TableRowBlock) (result string, err error) {
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

//RequestRowValueByHeader: make a request to retrieve the config table and return the value of the corresponding header
func RequestRowValueByHeader(client *notionapi.Client, pageid string, header string) (result string, err error) {
	tableRow, err := RequestTableRowBlockByHeader(client, pageid, header)
	if err != nil {
		return "", err
	}
	return GetRowValue(tableRow)
}

//GetTargetUrlFromConfig: return the value of the cell specifying the target url/ip
func GetTargetUrlFromConfig(children notionapi.Blocks) (targetUrl string, err error) {
	tableRow, err := GetTableRowBlockbyHeader(children, TARGET)
	if err != nil {
		return "", err
	}
	return GetRowValue(tableRow)
}

//RequestTargetUrlFromConfig: return the value of the cell specifying the target url/ip by making a request
func RequestTargetUrlFromConfig(client *notionapi.Client, pageid string) (targetUrl string, err error) {
	return RequestRowValueByHeader(client, pageid, TARGET)
}

//GetTargetUrlFromConfig: return the value of the cell specifying the target url/ip
func GetPortFromConfig(children notionapi.Blocks) (port string, err error) {
	tableRow, err := GetTableRowBlockbyHeader(children, PORT)
	if err != nil {
		return "", err
	}
	return GetRowValue(tableRow)
}

//RequestPortFromConfig: return the value of the cell specifying the port by making a request
func RequestPortFromConfig(client *notionapi.Client, pageid string) (port string, err error) {

	return RequestRowValueByHeader(client, pageid, PORT)
}

//GetShellFromConfig: return the value of the cell specifying the shell
func GetShellFromConfig(children notionapi.Blocks) (shell string, err error) {
	tableRow, err := GetTableRowBlockbyHeader(children, SHELL)
	if err != nil {
		return "", err
	}
	return GetRowValue(tableRow)
}

//RequestShellFromConfig: return the value of the cell specifying the shell by making a request
func RequestShellFromConfig(client *notionapi.Client, pageid string) (shell string, err error) {

	return RequestRowValueByHeader(client, pageid, SHELL)
}

//AppendCodeBlock: Add a code block ("shell") at the end of the page
func AppendCodeBlock(client *notionapi.Client, pageid string, content string) (resp *notionapi.AppendBlockChildrenResponse, err error) {
	request := notionapi.AppendBlockChildrenRequest{
		Children: []notionapi.Block{
			&notionapi.CodeBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeCode,
				},
				Code: notionapi.Code{
					Language: "shell",
					RichText: []notionapi.RichText{
						{
							Type: notionapi.ObjectTypeText,
							Text: notionapi.Text{Content: content},
						},
					},
				},
			},
		},
	}
	resp, err = client.Block.AppendChildren(context.Background(), notionapi.BlockID(pageid), &request)
	return resp, err
}

//AppendEmbedBlock: Add a embed block with specified url
func AppendEmbedBlock(client *notionapi.Client, pageid string, url string) (resp *notionapi.AppendBlockChildrenResponse, err error) {
	request := notionapi.AppendBlockChildrenRequest{
		Children: []notionapi.Block{
			&notionapi.EmbedBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeEmbed,
				},
				Embed: notionapi.Embed{
					URL: url,
				},
			},
		},
	}
	resp, err = client.Block.AppendChildren(context.Background(), notionapi.BlockID(pageid), &request)
	return resp, err
}

//AppendEmbedBlock: Add a embed block with specified url
func AppendHeadingBlock(client *notionapi.Client, pageid string, content string) (resp *notionapi.AppendBlockChildrenResponse, err error) {
	//heading content
	var richs []notionapi.RichText
	rich := notionapi.RichText{
		Type: notionapi.ObjectTypeText,
		Text: notionapi.Text{
			Content: content,
		},
		Annotations: &notionapi.Annotations{
			Bold:   false,
			Italic: true,
			Code:   true,
			Color:  "green",
		},
	}
	richs = append(richs, rich)

	request := notionapi.AppendBlockChildrenRequest{
		Children: []notionapi.Block{
			&notionapi.Heading3Block{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeHeading3,
				},
				Heading3: notionapi.Heading{
					RichText: richs,
				},
			},
		},
	}
	resp, err = client.Block.AppendChildren(context.Background(), notionapi.BlockID(pageid), &request)
	return resp, err
}
