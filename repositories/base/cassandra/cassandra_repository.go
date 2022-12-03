package cassandra

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/novabankapp/common.data/repositories/base"
	"github.com/novabankapp/common.data/utils"
	"log"
	"time"

	"github.com/fatih/structs"
	"github.com/gocql/gocql"
	base2 "github.com/novabankapp/common.data/domain/base"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type CassandraRepository[E base2.NoSqlEntity] struct {
	session   *gocqlx.Session
	tableName string
	timeout   time.Duration
}

const (
	COLUMN  = "column"
	COMPARE = "compare"
	VALUE   = "value"
)

func NewCassandraRepository[E base2.NoSqlEntity](
	session *gocqlx.Session,
	tableName string, timeout time.Duration) base.NoSqlRepository[E] {
	return &CassandraRepository[E]{
		session:   session,
		tableName: tableName,
		timeout:   timeout,
	}
}
func (rep *CassandraRepository[E]) GetById(ctx context.Context, id string) (*E, error) {

	ctx, cancel := context.WithTimeout(ctx, rep.timeout)
	defer cancel()

	var result []E
	getUser := qb.Select(rep.tableName).
		Where(qb.EqLit("id", id)).
		Query(*rep.session).
		WithContext(ctx)

	err := getUser.Select(&result)
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, errors.New("record not found")
	}
	return &result[0], nil
}
func (rep *CassandraRepository[E]) Create(ctx context.Context, entity E) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, rep.timeout)
	defer cancel()
	fields := structs.Fields(entity)
	for _, field := range fields {
		if field.Name() == "ID" {
			field.Set(uuid.New().String())
		} else {
			continue
		}
	}
	columns := structs.Names(entity)
	insert := qb.Insert(rep.tableName).
		Columns(columns...).
		Query(*rep.session).
		WithContext(ctx)
	insert.BindStruct(entity)
	if err := insert.ExecRelease(); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
func (rep *CassandraRepository[E]) Update(ctx context.Context, entity E, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, rep.timeout)

	defer cancel()

	columns := utils.Map(structs.Names(entity), utils.ToSnakeCase)
	updateUser := qb.Update(rep.tableName).
		Set(columns...).
		Where(qb.EqLit("ID", id)).
		Query(*rep.session).
		SerialConsistency(gocql.Serial).WithContext(ctx)

	updateUser.BindStruct(entity)

	applied, err := updateUser.ExecCASRelease()
	if err != nil {
		return false, err
	}

	return applied, nil
}
func (rep *CassandraRepository[E]) Delete(ctx context.Context, id string) (bool, error) {

	ctx, cancel := context.WithTimeout(ctx, rep.timeout)
	defer cancel()

	ent, error := rep.GetById(ctx, id)

	if error != nil {
		return false, error
	}
	delete := qb.Delete(rep.tableName).
		Where(qb.EqLit("id", id)).
		Query(*rep.session).
		SerialConsistency(gocql.Serial).WithContext(ctx)

	delete.BindStruct(ent)
	applied, err := delete.ExecCASRelease()
	if err != nil {
		return false, err
	}

	return applied, nil
}

func (rep *CassandraRepository[E]) Get(ctx context.Context,
	page []byte, pageSize int, queries []map[string]string, orderBy string) (*[]E, []byte, error) {

	ctx, cancel := context.WithTimeout(ctx, rep.timeout)
	defer cancel()
	wheres := make([]qb.Cmp, len(queries))
	for query := range queries {
		m := queries[query]
		column := m[COLUMN]
		compare := m[COMPARE]
		value := m[VALUE]
		var where qb.Cmp
		switch compare {
		case "<":
			where = qb.LtLit(column, value)
		case "<=":
			where = qb.LtOrEqLit(column, value)
		case ">":
			where = qb.GtLit(column, value)
		case ">=":
			where = qb.GtOrEqLit(column, value)
		case "=":
			where = qb.EqLit(column, value)
		default:
			where = qb.EqLit(column, value)

		}
		wheres = append(wheres, where)

	}
	var results []E

	get := qb.Select(rep.tableName).
		OrderBy(orderBy, qb.DESC)
	for i := range wheres {
		get = get.Where(wheres[i])
	}
	itr := get.
		Query(*rep.session).
		PageSize(pageSize).
		PageState(page).
		WithContext(ctx).
		Iter()

	page = itr.PageState()

	err := itr.Select(&results)
	if err != nil {
		return nil, nil, err
	}
	return &results, page, nil
}
func (rep *CassandraRepository[E]) GetByCondition(ctx context.Context,
	queries []map[string]string) (*E, error) {

	ctx, cancel := context.WithTimeout(ctx, rep.timeout)
	defer cancel()
	wheres := make([]qb.Cmp, len(queries))
	for query := range queries {
		m := queries[query]
		column := m[COLUMN]
		compare := m[COMPARE]
		value := m[VALUE]
		var where qb.Cmp
		switch compare {
		case "<":
			where = qb.LtLit(column, value)
		case "<=":
			where = qb.LtOrEqLit(column, value)
		case ">":
			where = qb.GtLit(column, value)
		case ">=":
			where = qb.GtOrEqLit(column, value)
		case "=":
			where = qb.EqLit(column, value)
		default:
			where = qb.EqLit(column, value)

		}
		wheres = append(wheres, where)

	}
	var results []E
	get := qb.Select(rep.tableName)
	for i := range wheres {
		get = get.Where(wheres[i])
	}
	query := get.
		Query(*rep.session).
		WithContext(ctx)

	err := query.Select(&results)
	if err != nil {
		return nil, err
	}
	return &results[0], nil
}
