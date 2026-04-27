package sysdept

import (
	"context"
	"system/model"

	"github.com/yitter/idgenerator-go/idgen"
)

// Add 增
func (m *SysDeptImpl) Add(ctx context.Context, args *model.SysDeptDto) (int64, error) {
	v := args.SysDept
	v.Id = idgen.NextId()
	v.DeleteTime = 0
	return m.engine.Context(ctx).InsertOne(&v)
}
