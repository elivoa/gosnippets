package pg

import (
	"sync"
)

// TODO need a graceful shutdown prorgress.

// TODO db config use config
// var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", "bogao", "", "localhost", "5432", "todo")
var mockConfig = &PGDatabaseConfig{
	Key:         "default",
	Host:        "localhost",
	Port:        "5432",
	Database:    "todo",
	Username:    "bogao",
	Password:    "",
	MaxPoolSize: new(int),
}

// named connect pool
var dbpoolcache sync.Map

// default driver instance
var defaultDriver *DatabaseDriver
var once sync.Once

// for named driver. for muitl-connection
func GetDriver(key string) (*DatabaseDriver, error) {
	if v, ok := dbpoolcache.Load(key); ok {
		if driver, ok := v.(*DatabaseDriver); ok {
			return driver, nil
		}
	}
	// driver := &DatabaseDriver{Config: mockConfig}
	driver, err := NewDriver(mockConfig)
	if err != nil {
		return nil, err
	}
	dbpoolcache.Store(key, driver)
	return driver, nil
}

func GetDefaultDriver() (driver *DatabaseDriver, erro error) {
	once.Do(func() {
		d, err := GetDriver("default") // get default driver
		if err != nil {
			erro = err
		}
		defaultDriver = d
		driver = d
	})
	return defaultDriver, nil
}

func GetDefaultDriverPanic() *DatabaseDriver {
	driver, err := GetDefaultDriver()
	if err != nil {
		panic(err)
	}
	return driver
}

func CloseDriver(key string) {
	if v, ok := dbpoolcache.Load(key); ok {
		if driver, ok := v.(*DatabaseDriver); ok {
			driver.pool.Close()
		}
	}
}

// func TestConnExec(t *testing.T) {
// 	t.Parallel()

// 	pool, err := pgxpool.Connect(context.Background(), os.Getenv("PGX_TEST_DATABASE"))
// 	require.NoError(t, err)
// 	defer pool.Close()

// 	c, err := pool.Acquire(context.Background())
// 	require.NoError(t, err)
// 	defer c.Release()

// 	testExec(t, c)
// }
