package middleware

import (
	"strings"

	"pkg/apicrypt"
	"pkg/conf"

	"github.com/gofiber/fiber/v2"
)

// ApiCrypto 处理 X-Shh-Encrypted 约定：解密请求体、标记响应需加密（实际加密在 pkg/resp 写入时完成）
func ApiCrypto() fiber.Handler {
	return func(c *fiber.Ctx) error {
		passphrase := strings.TrimSpace(conf.App.EncryptKey)
		if shouldSkipApiCrypto(c.Path()) {
			return c.Next()
		}

		clientWants := c.Get(apicrypt.HeaderName) == apicrypt.HeaderValue
		serverOn := passphrase != ""
		needEncResp := serverOn && clientWants
		if needEncResp {
			c.Locals(apicrypt.LocalsEncResponse, true)
		}

		if !serverOn {
			return c.Next()
		}

		if clientWants {
			ct := strings.ToLower(string(c.Request().Header.ContentType()))
			if strings.Contains(ct, "multipart/form-data") {
				return c.Next()
			}
			if hasJSONBody(c) && len(c.Body()) > 0 {
				plain, err := apicrypt.DecryptEnvelopeToPlain(c.Body(), passphrase)
				if err != nil {
					return fiber.NewError(fiber.StatusBadRequest, "invalid encrypted body")
				}
				c.Request().SetBodyRaw(plain)
				c.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
			}
		}

		return c.Next()
	}
}

func hasJSONBody(c *fiber.Ctx) bool {
	switch c.Method() {
	case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
	default:
		return false
	}
	ct := strings.ToLower(string(c.Request().Header.ContentType()))
	return strings.Contains(ct, fiber.MIMEApplicationJSON) || strings.Contains(ct, "application/json")
}

func shouldSkipApiCrypto(path string) bool {
	if strings.Contains(path, "/sys_websocket") || strings.Contains(path, "/sys_sse") {
		return true
	}
	if strings.Contains(path, "/dev_process/run_stream") {
		return true
	}
	if strings.Contains(path, "/sys_im/chat") {
		return true
	}
	return false
}
