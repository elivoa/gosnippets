package pg_test

import (
	"context"
	"fmt"
	"infra/driver/pg"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

// TODO config this.
var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", "bogao", "", "localhost", "5432", "todo")

func TestConnExec(t *testing.T) {

	pool, err := pgxpool.Connect(context.Background(), dbURL) // os.Getenv("PGX_TEST_DATABASE")
	require.NoError(t, err)
	defer pool.Close()

	c, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer c.Release()

	// testExec(t, c)
}

func TestConnExecPool(t *testing.T) {
	pgd, err := pg.GetDefaultDriver()
	if err != nil {
		t.Error(err)
	}
	defer pgd.ClosePool()

	pool, err := pgd.GetPool()
	if err != nil {
		t.Error(err)
	}

	// pool, err := pgxpool.Connect(context.Background(), dbURL) // os.Getenv("PGX_TEST_DATABASE")
	// require.NoError(t, err)
	// defer pool.Close()

	c, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer c.Release()

	// testExec(t, c)
}

func TestConnExecPoolQuick(t *testing.T) {
	pgd := pg.GetDefaultDriverPanic()
	defer pgd.ClosePool()

	pool := pgd.GetPoolPanic()

	c, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer c.Release()

	// testExec(t, c)
}
