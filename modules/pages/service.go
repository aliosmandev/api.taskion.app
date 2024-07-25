package pages

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

func getPages(c *fiber.Ctx) error {
	var postUrl string = "https://api.notion.com/v1/search"

	payload := NotionSearchPayload{
		Sort: Sort{
			Direction: "ascending",
			Timestamp: "last_edited_time",
		},
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return c.JSON(500, "error")
	}

	buffer := bytes.NewBuffer(payloadJson)

	r, err := http.NewRequest("POST", postUrl, buffer)
	if err != nil {
		return c.JSON(500, "error")
	}

	Authorization := c.GetReqHeaders()["Authorization"]
	accessToken := Authorization[len(Authorization)-1][7:]

	log.Info("token geldi ", accessToken)

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
