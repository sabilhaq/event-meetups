package main

type config struct {
	Port    string        `cfg:"port" cfgDefault:"9186"`
	Storage storageConfig `cfg:"storage"`
}

type storageConfig struct {
	Type     storageType           `cfg:"type" cfgDefault:"mysql"`
	Memory   storageMemoryConfig   `cfg:"memory"`
	DynamoDB storageDynamoDBConfig `cfg:"dynamodb"`
	MySQL    storageMySQLConfig    `cfg:"mysql"`
	SMTP     storageSMTPConfig     `cfg:"smtp"`
}

type storageMemoryConfig struct {
	MonsterDataPath string `cfg:"monster_data_path" cfgDefault:"../../deploy/local/run/rest-memory/data.json"`
	UserDataPath    string `cfg:"user_data_path" cfgDefault:"../../deploy/local/run/rest-memory/users.json"`
	EventDataPath   string `cfg:"event_data_path" cfgDefault:"../../deploy/local/run/rest-memory/events.json"`
	VenueDataPath   string `cfg:"venue_data_path" cfgDefault:"../../deploy/local/run/rest-memory/venues.json"`
}

type storageDynamoDBConfig struct {
	LocalstackEndpoint string `cfg:"localstack_endpoint"`
	BattleTableName    string `cfg:"battle_table_name"`
	GameTableName      string `cfg:"game_table_name"`
	MonsterTableName   string `cfg:"monster_table_name"`
}

type storageMySQLConfig struct {
	SQLDSN string `cfg:"sql_dsn" cfgDefault:"root:password@tcp(localhost:3306)/db_eventmeetup?timeout=5s"`
}

type storageSMTPConfig struct {
	Host      string `cfg:"smtp_host" cfgDefault:"localhost"`
	Port      string `cfg:"smtp_port" cfgDefault:"1025"` // Default port for Mailpit
	FromEmail string `cfg:"smtp_from_email" cfgDefault:"no-reply@yourdomain.com"`
}

type storageType string

const (
	storageTypeMemory   storageType = "memory"
	storageTypeMySQL    storageType = "mysql"
	storageTypeDynamoDB storageType = "dynamodb"
)
