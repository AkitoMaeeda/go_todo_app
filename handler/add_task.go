package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AkitoMaeeda/go_todo_app/entity"
	"github.com/AkitoMaeeda/go_todo_app/store"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type AddTask struct {

	//データベースを保存先として使用するため、TaskStoreにタスクを保存する必要がなくなった。
	//Store     *store.TaskStore

	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}

	//json型のhttp.requestをデコードし、bodyをbに代入
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//bのTitleが存在しているかどうか調査
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//タスクをTask構造体に登録
	t := &entity.Task{
		Title:  b.Title,
		Status: entity.TaskStatusTodo,
	}

	//登録したタスクを保存
	err := at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.id}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
