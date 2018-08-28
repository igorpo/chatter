package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/igorpo/chatter/chatroom"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// implements the handler interface for a template handler
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := chatroom.NewRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.Run()
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Server could not startup: ", err)
	}
}
