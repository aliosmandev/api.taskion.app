package blocks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type NotionSearchPayload struct {
	Sort Sort `json:"sort"`
}

type Sort struct {
	Direction string `json:"direction"`
	Timestamp string `json:"timestamp"`
}

func getBlocks(c *fiber.Ctx) error {

	var pageId string = c.Params("pageId")

	var getUrl string = "https://api.notion.com/v1/blocks/" + pageId + "/children?page_size=100"

	r, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return c.JSON(500, "error")
	}

	Authorization := c.GetReqHeaders()["Authorization"]
	accessToken := Authorization[len(Authorization)-1][7:]

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+accessToken)
	r.Header.Set("Notion-Version", "2022-02-22")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return c.JSON(500, "error")
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return c.Status(res.StatusCode).JSON(string(bodyBytes))
	}

	var responseBody map[string]interface{}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return c.JSON(500, "error")
	}

	return c.JSON(responseBody)
}

type AddTodoPayload struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
}

type RichText struct {
	Annotations Annotations `json:"annotations"`
	Href        *string     `json:"href"`
	PlainText   string      `json:"plain_text"`
	Text        struct {
		Content string      `json:"content"`
		Link    interface{} `json:"link"`
	} `json:"text"`
	Type string `json:"type"`
}

// ToDo represents the structure of the to_do field
type ToDo struct {
	Checked  bool       `json:"checked"`
	Color    string     `json:"color"`
	RichText []RichText `json:"rich_text"`
}

// Block represents the structure of each block
type Block struct {
	Object string `json:"object"`
	ToDo   ToDo   `json:"to_do"`
	Type   string `json:"type"`
}

type NotionBlock struct {
	Children []Block `json:"children"`
}

func AddTodo(c *fiber.Ctx) error {
	var pageId string = c.Params("pageId")
	var postUrl string = "https://api.notion.com/v1/blocks/" + pageId + "/children"

	payload := new(AddTodoPayload)

	if err := c.BodyParser(payload); err != nil {
		log.Info("error = ", err)
		return c.SendStatus(200)
	}

	var requestPayload NotionBlock

	requestPayload.Children = append(requestPayload.Children, Block{
		Object: "block",
		Type:   "to_do",
		ToDo: ToDo{
			Checked: payload.Checked,
			Color:   "default",
			RichText: []RichText{
				{
					Annotations: Annotations{
						Bold:          false,
						Code:          false,
						Color:         "default",
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
					},
					Text: struct {
						Content string      `json:"content"`
						Link    interface{} `json:"link"`
					}{
						Content: payload.Text,
					},
					PlainText: payload.Text,
					Type:      "text",
				},
			},
		},
	})

	payloadJson, err := json.Marshal(requestPayload)
	if err != nil {
		return c.JSON(500, "error")
	}

	buffer := bytes.NewBuffer(payloadJson)

	r, err := http.NewRequest("PATCH", postUrl, buffer)
	if err != nil {
		return c.JSON(500, "error")
	}

	Authorization := c.GetReqHeaders()["Authorization"]
	accessToken := Authorization[len(Authorization)-1][7:]

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+accessToken)
	r.Header.Set("Notion-Version", "2022-02-22")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return c.JSON(500, "error")
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return c.Status(res.StatusCode).JSON(string(bodyBytes))
	}

	var responseBody map[string]interface{}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return c.JSON(500, "error")
	}

	return c.JSON(responseBody)
}
