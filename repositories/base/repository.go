package base

import (
	"context"
	domainbase "github.com/novabankapp/common.data/domain/base"
)

type RdbmsRepository[E domainbase.Entity] interface {
	Create(ctx context.Context, entity E) (*E, error)
	Update(ctx context.Context, entity E, id uint) (bool, error)
	Delete(ctx context.Context, id uint) (bool, error)
	GetById(ctx context.Context, ID string) (*E, error)
	Get(ctx context.Context, page int, pageSize int, query *E, orderBy string) (*[]E, error)
}

type NoSqlRepository[E domainbase.NoSqlEntity] interface {
	GetById(ctx context.Context, id string) (*E, error)
	Create(ctx context.Context, entity E) (bool, error)
	Update(ctx context.Context, entity E, id string) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, page []byte, pageSize int, queries []map[string]string, orderBy string) (*[]E, error)
	GetByCondition(ctx context.Context, queries []map[string]string) (*E, error)
}
