package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AkitoMaeeda/go_todo_app/config"
	_ "github.com/jmoiron/sqlx"
	"golang.org/x/mobile/exp/sprite/clock"
)

// この辺が何で必要なのか今のところ不明。
type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, quely string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, quely string, qrgs ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, quely string, args ...any) error
}

var (
	//上で宣言した構造体が実際の型シグネチャ(名前とか引数とか)と一致しているのか確認する
	/*細かくすると、本来なら　var b *sqlx.DB = nil , var _　Beginner = b
	↓これ（ _ ）変数名ね。*/
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

// 永続化をrepositry型で実行するためらしい。時間を取得するオブジェクトを取得
type Repository struct {
	Clocker clock.Clocker
}

// config.Config型を使ってDBへ接続する
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {

	//DBへ接続する（接続確認は行っていない）
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime = true",
			cfg.DBUser, cfg.DBPassword,
			cfg.DBHost, cfg.DBPort,
			cfg.DBName,
		),
	)

	if err != nil {
		return nil, nil, err
	}

	//2秒後にcancelされるように設定
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	//接続の生存確認
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}

	//sql.DBに対する新しいsql.DBラッパーを作成
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil

}
