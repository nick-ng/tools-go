package jira

import "fmt"

type Mark struct {
	Type  string            `json:"type"`
	Attrs map[string]string `json:"attrs"`
}

type Content struct {
	Type    string    `json:"type"`
	Content []Content `json:"content"`
	Text    string    `json:"text"`
	Marks   []Mark    `json:"marks"`
}

// @todo(nick-ng): figure out how to handle bulletlists
func ContentToString(contents []Content, input string) string {
	temp := input

	for _, content := range contents {
		switch content.Type {
		case "text":
			temp2 := textToString(content)
			temp = fmt.Sprintf("%s%s", temp, temp2)
		case "paragraph":
			{
				temp2 := ContentToString(content.Content, temp)
				temp = fmt.Sprintf("%s\n\n", temp2)
			}
		default:
			{
				// fmt.Println(content)
			}
		}
	}

	return temp
}

func textToString(content Content) string {
	for _, mark := range content.Marks {
		switch mark.Type {
		case "link":
			{
				if len(mark.Attrs["href"]) > 0 && content.Text != mark.Attrs["href"] {
					return fmt.Sprintf("[%s](\x1b[4m\x1b[36m%s\x1b[0m)", content.Text, mark.Attrs["href"])
				}

				return fmt.Sprintf("\x1b[4m\x1b[36m%s\x1b[0m", content.Text)
			}
		}
	}

	return content.Text
}
