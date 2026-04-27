package sysauth

import (
	"os"
	"path/filepath"
	"pkg/auth"
	"pkg/conf"
	"pkg/constx"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"
	"strconv"
	"strings"
	sysauth "system/impl/sys_auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	avatarMaxSize   = 5 << 20 // 5MB
	avatarFormField = "file"
)

var allowedAvatarTypes = map[string]bool{
	"image/jpeg": true, "image/jpg": true, "image/png": true,
	"image/gif": true, "image/webp": true,
}

var allowedAvatarExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

// @Summary		系统认证_上传个人头像
// @Description	当前登录用户上传头像图片，支持 jpg/png/gif/webp，最大 5MB
// @Tags			SysAuthApi
// @Security		ApiKeyAuth
// @Accept			multipart/form-data
// @Produce		json
// @Param			file	formData	file	true	"头像图片"
// @Success		200		{object}	r.MyResp{data=UploadAvatarResp}	"Success"
// @Failure		500		{object}	r.MyResp	"Failure"
// @Router			/sys_auth/avatar [post]
func UploadAvatar(c *fiber.Ctx) error {
	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}

	file, err := c.FormFile(avatarFormField)
	if err != nil {
		return errs.Args(err)
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedAvatarTypes[contentType] {
		return errs.New("仅支持 jpg、png、gif、webp 格式图片")
	}
	if file.Size > avatarMaxSize {
		return errs.New("图片大小不能超过 5MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	if !allowedAvatarExts[ext] {
		return errs.New("仅支持 jpg、png、gif、webp 扩展名")
	}
	fileName := strconv.FormatInt(loginId, 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	dir := filepath.Join(conf.File.Local, "avatar")
	if err = os.MkdirAll(dir, 0755); err != nil {
		return errs.Sys(err)
	}
	savePath := filepath.Join(dir, fileName)
	if err = c.SaveFile(file, savePath); err != nil {
		return errs.Sys(err)
	}

	// 存库使用相对路径，便于前端拼接 baseUrl
	avatarPath := "avatar/" + fileName
	sysAuthImpl := sysauth.Impl()
	if err = sysAuthImpl.UpdateAvatar(c.Context(), datacv.IntToStr(loginId), avatarPath); err != nil {
		return err
	}
	// 同步更新 session，供 IM 聊天室等使用
	_ = auth.SetSess(c.Context(), datacv.IntToStr(loginId), map[string]any{"avatar": avatarPath})

	return r.Resp(c, UploadAvatarResp{Avatar: conf.FileUrl(avatarPath)})
}
