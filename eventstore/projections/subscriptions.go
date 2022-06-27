package projections

type Subscriptions struct {
	PoolSize                     int    `mapstructure:"poolSize" validate:"required,gte=0"`
	OrderPrefix                  string `mapstructure:"orderPrefix" validate:"required,gte=0"`
	CassandraProjectionGroupName string `mapstructure:"cassandraProjectionGroupName" validate:"required,gte=0"`
	ElasticProjectionGroupName   string `mapstructure:"elasticProjectionGroupName" validate:"required,gte=0"`
}
