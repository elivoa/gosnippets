package pg

import "fmt"

// TODO load this with config file.

type PGDatabaseConfig struct {
	// each connection has a key.
	// "" or "default" as the default db config.
	Key string

	// basic connection
	Host     string
	Port     string
	Database string
	Username string
	Password string

	// connection pool related configs.
	MaxPoolSize *int

	// TODO default timeout & configured timeout.

	// todo more configs.
}

func (p *PGDatabaseConfig) ToConnectURL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		p.Username, p.Password, p.Host, p.Port, p.Database,
	)
}

// ---------------------------------------------------------------------------
// Config Suit, support more than one database connection in one applicaion.
// ---------------------------------------------------------------------------

type DBConfigSuit struct {
	DBConn []*PGDatabaseConfig
}

// TODO lazy load config?

func GetKeyedDBConfig(key string) *PGDatabaseConfig {
	return nil // TODO implement GetKeyedDBConfig
}

// ------------------------------- mock data -------------------------------

// var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", "bogao", "", "localhost", "5432", "todo")
