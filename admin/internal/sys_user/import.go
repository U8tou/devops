package sysuser

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"pkg/constx"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/excel"
	"pkg/validator"
	sysuser "system/impl/sys_user"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// 表头列名（与 UserExcelHeaders 一致，用于从 map 取值）
const (
	colUserName  = "用户账号"
	colNickName  = "用户昵称"
	colUserType  = "用户类型"
	colEmail     = "邮箱"
	colPhoneArea = "电话区号"
	colPhone     = "电话号码"
	colSex       = "性别"
	colStatus    = "状态"
	colAddress   = "联系地址"
	colRemark    = "备注"
	colPassword  = "密码"

	importMaxSize = 10 << 20 // 10MB，防止大文件 DoS
)

// @Summary		用户导入
// @Description	上传 Excel 按模板导入用户（multipart/form-data）。file: 文件；updateExisting: 是否更新已存在的用户（true/1 为是）\n权限标识：sys:user:add
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			multipart/form-data
// @Produce		json
// @Param			file				formData		file	true	"用户导入模板填写的 xlsx 文件"
// @Param			updateExisting		formData		bool	false	"是否更新已经存在的用户数据"
// @Success		200		{object}	r.MyResp{data=ImportResp}	"Success"
// @Router			/sys_user/import [post]
func Import(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	if file.Size > importMaxSize || file.Size <= 0 {
		return r.Error(c, errs.Args(errors.New("导入文件大小不能超过 10MB")))
	}
	fh, err := file.Open()
	if err != nil {
		return err
	}
	defer fh.Close()

	// 是否更新已存在的用户（表单字段，勾选时多为 "true" 或 "1"）
	updateExisting := parseBoolForm(c.FormValue("updateExisting"))

	_, rows, err := excel.ReadSheetFromReader(fh)
	if err != nil {
		return err
	}

	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}
	tp := time.Now().Unix()

	impl := sysuser.Impl()
	var successCount int
	var failList []ImportFail

	for i, row := range rows {
		excelRowNum := i + 2 // 表头为第 1 行
		args, errMsg := parseAndValidateRow(row)
		if errMsg != "" {
			failList = append(failList, ImportFail{Row: excelRowNum, Reason: errMsg})
			continue
		}
		args.Avatar = ""
		args.UpdateTime = tp
		args.UpdateBy = loginId

		existing, err := impl.GetByUserName(c.Context(), args.UserName)
		if err != nil {
			failList = append(failList, ImportFail{Row: excelRowNum, Reason: err.Error()})
			continue
		}
		if existing != nil {
			// 已存在：仅当勾选“更新已有数据”时更新，否则跳过
			if !updateExisting {
				failList = append(failList, ImportFail{Row: excelRowNum, Reason: "用户已存在，未勾选更新已有数据时跳过"})
				continue
			}
			args.Id = existing.Id
			_, err = impl.Edit(c.Context(), &args)
		} else {
			// 不存在：新增（密码必填）
			if args.Password == "" {
				failList = append(failList, ImportFail{Row: excelRowNum, Reason: "新增用户时密码不能为空"})
				continue
			}
			args.CreateTime = tp
			args.CreateBy = loginId
			_, err = impl.Add(c.Context(), &args)
		}
		if err != nil {
			failList = append(failList, ImportFail{Row: excelRowNum, Reason: err.Error()})
			continue
		}
		successCount++
	}

	return r.Resp(c, ImportResp{SuccessCount: successCount, FailList: failList})
}

func parseAndValidateRow(row map[string]string) (model.SysUserDto, string) {
	var dto model.SysUserDto
	get := func(key string) string { return strings.TrimSpace(row[key]) }

	userName := get(colUserName)
	if userName == "" {
		return dto, "用户账号不能为空"
	}
	if len(userName) > 30 {
		return dto, "用户账号长度不能超过30"
	}
	dto.UserName = userName

	nickName := get(colNickName)
	if nickName == "" {
		return dto, "用户昵称不能为空"
	}
	dto.NickName = nickName

	userType := get(colUserType)
	if userType == "" {
		userType = "00"
	}
	dto.UserType = userType

	email := get(colEmail)
	if len(email) > 50 {
		return dto, "邮箱长度不能超过50"
	}
	if email != "" {
		if err := validator.DefaultValidator.Var(email, "email"); err != nil {
			return dto, "邮箱格式不正确"
		}
	}
	dto.Email = email

	phoneArea := get(colPhoneArea)
	if phoneArea == "" {
		phoneArea = "+86"
	}
	dto.PhoneArea = phoneArea

	phone := get(colPhone)
	if phone == "" {
		return dto, "电话号码不能为空"
	}
	dto.Phone = phone

	sex, err := parseSmallInt(get(colSex), 1, 1, 3)
	if err != nil {
		return dto, "性别必须为 1(男)/2(女)/3(未知)"
	}
	dto.Sex = int8(sex)

	status, err := parseSmallInt(get(colStatus), 1, 1, 2)
	if err != nil {
		return dto, "状态必须为 1(正常)/2(停用)"
	}
	dto.Status = int8(status)

	dto.Address = get(colAddress)
	dto.Remark = get(colRemark)

	pwd := get(colPassword)
	// 密码：更新时可为空（表示不修改）；新增时在调用 Add 前校验必填
	if pwd != "" {
		if len(pwd) < 6 {
			return dto, "密码不能少于6位"
		}
		if len(pwd) > 100 {
			return dto, "密码不能超过100位"
		}
	}
	dto.Password = pwd

	return dto, ""
}

// parseBoolForm 解析表单中的布尔值（如复选框），"true"/"1"/"on"/"yes" 为 true，其余为 false
func parseBoolForm(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "true" || s == "1" || s == "on" || s == "yes"
}

// parseSmallInt 解析整数，支持 "1" 或 "1.0"；若为空则返回 defaultVal；必须在 [min,max] 内
func parseSmallInt(s string, defaultVal, min, max int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		// Excel 可能为 1.0
		f, e := strconv.ParseFloat(s, 64)
		if e != nil {
			return 0, e
		}
		n = int(f)
	}
	if n < min || n > max {
		return 0, strconv.ErrSyntax
	}
	return n, nil
}
