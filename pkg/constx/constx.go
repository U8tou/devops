package constx

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetLoginId(c *fiber.Ctx, loginId string) {
	c.Locals("userID", loginId)
}

func GetLoginId(c *fiber.Ctx) (int64, error) {
	userIDAny := c.Locals("userID")
	userID, ok := userIDAny.(string)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return strconv.ParseInt(userID, 10, 64)

}
