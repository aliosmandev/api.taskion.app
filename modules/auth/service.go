package auth

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Authorize(c *fiber.Ctx) error {
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

	fmt.Println(base.String())

	return c.Redirect(base.String())
}
