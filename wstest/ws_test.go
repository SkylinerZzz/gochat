package wstest

import (
	"github.com/gorilla/websocket"
	"gochat/pkg/service"
	"gochat/util"
	"html/template"
	"net/http"
	"testing"
)

func TestWs(t *testing.T) {
	util.Init("../config")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tp, err := template.ParseFiles("index.html")
		if err != nil {
			t.Fatal(err)
		}
		tp.Execute(w, nil)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, _ := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		entry := service.NewEntry()
		entry.Exec(ws, "1", "1")
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		t.Fatal(err)
	}
}
