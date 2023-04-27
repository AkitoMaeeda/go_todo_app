
package store

import (
	"context"
	"database/sql"
	"entity"
	"github.com/jmoiron/sqlx"
)

//select文で全部のタスクを取得するで！！
func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT id, title, status, created, modified FROM task;`

	//sqlx.Selectcontextは複数のレコードを取得し、各レコードを一つの構造体に
	//代入したスライスを返す。ここだと第3引数で取得したレコードをtask構造体に代入する...?

	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

//タスクの追加
func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	//idを指定していないということは、自動的に追加される？
	sql := `INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?)`
	//sql文の実行
	result, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.Created, t.Modified)

	if err != nil {
		return err
	}

	//オートインクリメントされたidの取得(schema.sqlで設定済み)
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}
	//idの保存
	t.ID = entity.TaskID(id)
	return nil
}
