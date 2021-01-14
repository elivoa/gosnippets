package staffdb

import (
	"context"
	"fmt"
	"infra/driver/pg"
	"mag/apps/kerrigan/countermgr/cmmodels"
	"strings"

	_ "github.com/doug-martin/goqu/v8/dialect/postgres"
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/pgxscan"
)

func GetUser(uid string) (*cmmodels.Staff, error) {
	// 标准流程
	// TODO quick流程
	pgd, err := pg.GetDefaultDriver()
	if err != nil {
		return nil, err
	}

	pool, err := pgd.GetPool()
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	fmt.Println("..", pool)

	rows, err := conn.Query(context.Background(), "SELECT * FROM tasks WHERE id = $1", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("??????????????????????", rows)

	// var pgxcon interface{} = conn
	// c := pgxcon.(*stdlib.Conn)
	// fmt.Println("..", c)

	// d := pgxcon.(*sql.DB)

	// fmt.Println("..d ", d)

	dialect := goqu.Dialect("postgres")
	// dialect := goqu.New("postgres", c)

	// TODO auto process database name
	// sd := dialect.From("tasks").Where(goqu.Ex{"id": 1})

	// TODO auto process database name
	ppd := dialect.From("tasks").Prepared(true)
	sd := ppd.Where(goqu.Ex{"id": 1})
	sql, args, err := sd.ToSQL()
	if err != nil {
		return nil, err
	}

	fmt.Println("1213123123123123123???????????", sql, args)

	var model cmmodels.Staff
	sql = strings.ReplaceAll(sql, "?", "$1")

	fmt.Println("--0-0-0-0-0-", sql, args)

	if err := pgxscan.Get(context.Background(), pool, &model, sql, args...); err != nil {
		panic(err)
		return nil, err
	}
	fmt.Println(">>>>>>>>>>>>>>>", model)

	// success, err := sd.ScanStruct(&model)
	// if err != nil {
	// 	return nil, err
	// }

	// ds := dialect.From("todo").Where(goqu.Ex{"id": 10})
	// sql, args, err := ds.ToSQL()
	// if err != nil {
	// 	fmt.Println("An error occurred while generating the SQL", err.Error())
	// } else {
	// 	fmt.Println(sql, args)
	// }
	return nil, err
}
