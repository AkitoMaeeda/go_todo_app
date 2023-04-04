package main

import (
	"net/http"

	"github.com/AkitoMaeeda/go_todo_app/handler"
	"github.com/AkitoMaeeda/go_todo_app/store"
	"github.com/go-playground/validator/v10"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "application/json; charset = utf-8")

		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	//バリデーションのインスタンス作成
	v := validator.New()

	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServerHTTP)

	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux
}
