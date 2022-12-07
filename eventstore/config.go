package eventstore

// Config of es package.
type Config struct {
	SnapshotFrequency int64 `json:"snapshotFrequency" validate:"required,gte=0"`
}

type EventStoreConfig struct {
	ConnectionString string `mapstructure:"connectionString"`
}
