package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 33306
	//環境変数の取得。今回だとCIはgithubActionsに設定されている
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306
	}
	db, err := sql.Open(
		"mysql",
		fmt.Sprint("todo:todo@tcp(127.0,0,1:%d)/todo?parseTime = true", port),
	)

	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(
		func() { _ = db.Close() },
	)

	//sqlx.db型を取得する
	return sqlx.NewDb(db, "mysql")

}
