package sysauth

import (
	"context"
	"pkg/crypto"
	"pkg/errs"
	"system/model"
	"time"

	sysuser "system/impl/sys_user"

	"github.com/yitter/idgenerator-go/idgen"
)

func (m *SysAuthImpl) Register(ctx context.Context, userName string, password string) error {
	userImpl := sysuser.Impl()
	existing, err := userImpl.GetByUserName(ctx, userName)
	if err != nil {
		return errs.Sys(err)
	}
	if existing != nil {
		return errs.ERR_HAS_ACCOUNT
	}

	hashedPwd, err := crypto.HashPassword(password)
	if err != nil {
		return errs.Sys(err)
	}

	tp := time.Now().Unix()
	sess := m.engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return errs.Sys(err)
	}
	count, err := sess.Context(ctx).Count(&model.SysUser{})
	if err != nil {
		_ = sess.Rollback()
		return errs.Sys(err)
	}
	userType := model.UserTypeNormal
	if count == 0 {
		userType = model.UserTypeRoot
	}
	user := model.SysUser{
		Id:         idgen.NextId(),
		UserName:   userName,
		NickName:   userName,
		UserType:   userType,
		Email:      "",
		PhoneArea:  "+86",
		Phone:      "",
		Sex:        3,
		Avatar:     "",
		Password:   hashedPwd,
		Status:     1,
		Address:    "",
		Remark:     "",
		CreateTime: tp,
		CreateBy:   0,
		UpdateTime: tp,
		UpdateBy:   0,
		DeleteTime: 0,
	}
	_, err = sess.Context(ctx).InsertOne(&user)
	if err != nil {
		_ = sess.Rollback()
		return errs.Sys(err)
	}
	if err := sess.Commit(); err != nil {
		return errs.Sys(err)
	}
	return nil
}
