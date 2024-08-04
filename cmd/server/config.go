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
}

type storageMemoryConfig struct {
	MonsterDataPath string `cfg:"monster_data_path" cfgDefault:"../../deploy/local/run/rest-memory/data.json"`
	EventDataPath   string `cfg:"event_data_path" cfgDefault:"../../deploy/local/run/rest-memory/events.json"`
	UserDataPath    string `cfg:"user_data_path" cfgDefault:"../../deploy/local/run/rest-memory/users.json"`
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

type storageType string

const (
	storageTypeMemory   storageType = "memory"
	storageTypeMySQL    storageType = "mysql"
	storageTypeDynamoDB storageType = "dynamodb"
)
