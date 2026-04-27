package resp

import (
	"encoding/json"
	"errors"
	"log"
	"pkg/apicrypt"
	"pkg/conf"
	"pkg/errs"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// MyResp 统一响应结构
// @Description 统一响应结构
type MyResp struct {
	Code      int    `json:"code"`           // 状态码
	Msg       string `json:"msg"`            // 消息
	Data      any    `json:"data,omitempty"` // 数据
	Timestamp string `json:"timestamp"`      // 时间戳
}

func sendJSON(c *fiber.Ctx, httpStatus int, body *MyResp) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	key := strings.TrimSpace(conf.App.EncryptKey)
	if key != "" {
		if v, ok := c.Locals(apicrypt.LocalsEncResponse).(bool); ok && v {
			enc, err := apicrypt.EncryptToEnvelopeJSON(raw, key)
			if err != nil {
				return err
			}
			c.Set(apicrypt.HeaderName, apicrypt.HeaderValue)
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Status(httpStatus).Send(enc)
		}
	}
	return c.Status(httpStatus).JSON(body)
}

func Ok(c *fiber.Ctx) error {
	resp := &MyResp{
		Code:      200,
		Msg:       "success",
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}
	return sendJSON(c, fiber.StatusOK, resp)
}

func Resp(c *fiber.Ctx, data any) error {
	resp := &MyResp{
		Code:      200,
		Msg:       "success",
		Data:      data,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}
	return sendJSON(c, fiber.StatusOK, resp)
}

func Error(c *fiber.Ctx, err error) error {
	var resp MyResp
	resp.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	httpStatus := fiber.StatusOK
	// 自定义错误
	var errsErr *errs.Errs
	if errors.As(err, &errsErr) {
		resp.Code = errsErr.Code
		resp.Msg = errsErr.Msg
		resp.Data = errsErr.Data
		if errsErr.Code >= 400 && errsErr.Code < 600 {
			httpStatus = errsErr.Code
		}
		return sendJSON(c, httpStatus, &resp)
	}
	// 参数错误
	var fieldError validator.ValidationErrors
	if errors.As(err, &fieldError) {
		resp.Code = fiber.ErrBadRequest.Code
		resp.Msg = fiber.ErrBadRequest.Message
		resp.Data = fieldError.Error()
		return sendJSON(c, fiber.StatusBadRequest, &resp)
	}
	// 请求/响应异常
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		resp.Code = fiberErr.Code
		resp.Msg = fiberErr.Message
		return sendJSON(c, fiberErr.Code, &resp)
	}
	// 未知错误
	log.Println("❌ 系统异常:", err.Error())
	resp.Code = errs.ERR_SYS_DEFAULT.Code
	resp.Msg = errs.ERR_SYS_DEFAULT.Msg
	return sendJSON(c, fiber.StatusInternalServerError, &resp)
}
