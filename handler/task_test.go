package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AkitoMaeeda/go_todo_app/entity"
	"github.com/AkitoMaeeda/go_todo_app/store"
	"github.com/AkitoMaeeda/go_todo_app/testutil"
	"github.com/go-playground/validator/v10"
)

func TestAddTask(t *testing.T) {
	fmt.Println("t.parallel待機前")
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		repFile string
		want    want
	}{
		"ok": {
			repFile: "testdata/add_task/ok_req.json.golden",
			want: want{status: http.StasusOK,
				respFile: "testdata/add_task/bad_req_rsp.json.golden",
			},
		},

		"badRequest": {
			repFile: "testdata/add_task/bad_req.json.golden",
			want: want{status: http.StatusBadRequest,
				respFile: "testdata/add_task/bad_req_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		fmt.Println("t.Run待機前")
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}

}
