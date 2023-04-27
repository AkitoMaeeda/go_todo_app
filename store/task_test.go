package store

import (
	"context"
	"testing"
	"time"

	"github.com/AkitoMaeeda/go_todo_app/clock"
	"github.com/AkitoMaeeda/go_todo_app/entity"
	"github.com/AkitoMaeeda/go_todo_app/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}

	c := clock.FixedClocker{}

	wants := entity.Task{
		{
			Title: "want task 1", Status: "todo",
			Created: c.Now(), Modified: c, Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c, Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c, Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
			 VALUES
			 	(?, ?, ?, ?),(?, ?, ?, ?),(?, ?, ?, ?),(?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatalf(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTx(ctx, nil)

	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)

	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}

	var wantID int64 = 20

	okTask := &entity.Task{
		Titile:   "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	//戻り値の一つ目はモックdb接続、二つ目はモック設定を行うオブジェクト。
	db, mock, err := sqlmock.New()

	//ExpectExecのエラーもここで処理されるのかな...?
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	//トランザクションの開始を待機する。sqlライブラリ同様ExpectQuery()もある(その場合はWillReturnRows())
	//Withargsで？の引数を処理。
	//WillReturnResultの第一引数は「主キーの自動生成IDの値」、第二引数は本クエリによって影響を受けるカラムの数
	//今回だと、IDカラムが20のテーブルを出力する。
	mock.ExpectExec(

		`INSERT INTO task \(title, status, created, modified\) VALUES \(\?, \?, \?, \?)`,
	).WithArgs(okTask.Title, okTask.Status, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	//モックのDB接続sqlxを作成してる！？
	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
