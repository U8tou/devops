package syspost

import (
	"context"
	"system/model"

	"github.com/yitter/idgenerator-go/idgen"
)

// Add 增
func (m *SysPostImpl) Add(ctx context.Context, args *model.SysPostDto) (int64, error) {
	v := args.SysPost
	v.Id = idgen.NextId()
	v.DeleteTime = 0
	return m.engine.Context(ctx).InsertOne(&v)
}
