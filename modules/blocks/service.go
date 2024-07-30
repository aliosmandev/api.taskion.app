package blocks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

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

func createBlock(c *fiber.Ctx) error {
	var pageId string = c.Params("pageId")
	var postUrl string = "https://api.notion.com/v1/blocks/" + pageId + "/children"

	payload := new(TodoPayload)

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

type UpdateText struct {
	Content string `json:"content"`
}

type UpdateRichText struct {
	Text UpdateText `json:"text"`
}

type UpdateToDo struct {
	Checked  bool             `json:"checked"`
	RichText []UpdateRichText `json:"rich_text"`
}

type UpdateBlockPayload struct {
	ToDo UpdateToDo `json:"to_do"`
}

func updateBlock(c *fiber.Ctx) error {
	var blockId string = c.Params("blockId")

	var updateUrl string = "https://api.notion.com/v1/blocks/" + blockId

	payload := new(UpdateBlockPayload)

	if err := c.BodyParser(payload); err != nil {
		log.Info("error = ", err)
		return c.SendStatus(200)
	}

	var requestPayload UpdateBlockPayload

	requestPayload.ToDo.Checked = payload.ToDo.Checked
	requestPayload.ToDo.RichText = []UpdateRichText{
		{
			Text: UpdateText{
				Content: payload.ToDo.RichText[0].Text.Content,
			},
		},
	}

	payloadJson, err := json.Marshal(requestPayload)
	if err != nil {
		return c.JSON(500, "error")
	}

	buffer := bytes.NewBuffer(payloadJson)

	r, err := http.NewRequest("PATCH", updateUrl, buffer)
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

func deleteBlock(c *fiber.Ctx) error {
	var blockId string = c.Params("blockId")

	var deleteUrl string = "https://api.notion.com/v1/blocks/" + blockId

	r, err := http.NewRequest("DELETE", deleteUrl, nil)
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
