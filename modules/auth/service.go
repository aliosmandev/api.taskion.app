package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

type NotionTokenPayload struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}

type NotionTokenResponse struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	BotId         string `json:"bot_id"`
	WorkspaceId   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

func Authorize(c *fiber.Ctx) error {
	fmt.Println(os.Getenv("NOTION_AUTHORIZE_URL"))
	base, err := url.Parse(os.Getenv("NOTION_AUTHORIZE_URL"))
	if err != nil {
		return c.JSON(500, "error")
	}
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("owner", "user")
	params.Add("client_id", os.Getenv("NOTION_CLIENT_ID"))
	params.Add("redirect_uri", os.Getenv("NOTION_REDIRECT_URI"))

	base.RawQuery = params.Encode()

	return c.Redirect(base.String())
}

func Callback(c *fiber.Ctx) error {

	var code string = c.Query("code")

	var postUrl string = "https://api.notion.com/v1/oauth/token"

	body := []byte(`{
		"grant_type": "authorization_code",
		"code": "` + code + `",
		"redirect_uri": "` + os.Getenv("NOTION_REDIRECT_URI") + `"
	}`)

	r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(body))
	if err != nil {
		return c.JSON(500, "error")
	}

	var notionBearer string = base64.StdEncoding.EncodeToString([]byte(os.Getenv("NOTION_CLIENT_ID") + ":" + os.Getenv("NOTION_CLIENT_SECRET")))

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Basic "+notionBearer)
	r.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return c.JSON(500, "error")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return c.Status(res.StatusCode).JSON(string(bodyBytes))
	}

	var tokenResponse NotionTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return c.Status(500).JSON("error decoding response")
	}

	// Token alındı, şimdi token'ı kullanabilirsiniz
	return c.Status(200).JSON(tokenResponse)
}
