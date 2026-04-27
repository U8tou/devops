package sysauth

import (
	"context"
	"pkg/auth"
)

// Logout 登出
func (_ *SysAuthImpl) Logout(ctx context.Context, device string, token string) error {
	auth.OutToken(ctx, device, token)
	return nil
}
