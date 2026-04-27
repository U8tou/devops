package sysauth

import (
	"context"
	"pkg/auth"
	"system/model"
)

// RefreshToken 使用 RefreshToken 换取新的 AccessToken
func (_ *SysAuthImpl) RefreshToken(ctx context.Context, device string, refreshToken string) (*model.LoginResp, error) {
	newAccessToken, err := auth.RefreshToken(ctx, device, refreshToken)
	if err != nil {
		return nil, err
	}
	return &model.LoginResp{
		Token:        newAccessToken,
		RefreshToken: refreshToken, // 原 RefreshToken 保持不变
	}, nil
}
