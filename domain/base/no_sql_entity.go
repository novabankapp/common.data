package base

type NoSqlEntity interface {
	noSQLEntity
}

type noSQLEntity interface {
	IsNoSQLEntity() bool
}
