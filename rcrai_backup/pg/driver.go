package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO lazy? connect immediately.
var lazyConnect = true

type DatabaseDriver struct {
	Config      *PGDatabaseConfig
	initialized bool // TODO use once
	pool        *pgxpool.Pool
}

func NewDriver(config *PGDatabaseConfig) (*DatabaseDriver, error) {
	if nil == config {
		// TODO normal parameter check.
		return nil, fmt.Errorf("config Can't be nil") // TODO return an rcrai standard error with code?
	}
	driver := &DatabaseDriver{Config: config}

	if !lazyConnect {
		if err := driver.initpool(); err != nil {
			return nil, err
		}
	}

	return driver, nil
}

func (p *DatabaseDriver) initpool() (err error) {
	if !p.initialized {
		p.pool, err = pgxpool.Connect(context.Background(), p.Config.ToConnectURL())
		p.initialized = true
	}
	return
}

func (p *DatabaseDriver) GetPool() (pool *pgxpool.Pool, err error) {
	if !p.initialized {
		if lazyConnect {
			if err := p.initpool(); err != nil {
				return nil, err
			}
		}
	}
	return p.pool, nil
}

func (p *DatabaseDriver) GetPoolPanic() *pgxpool.Pool {
	if !p.initialized {
		if lazyConnect {
			if err := p.initpool(); err != nil {
				panic(err)
			}
		}
	}
	return p.pool
}

func (p *DatabaseDriver) ClosePool() {
	if p.pool != nil {
		p.pool.Close()
	}
}
