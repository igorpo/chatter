package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/igorpo/chatter/chatroom"
)

var (
	port = flag.Int("port", 8000, "The port the application will run on.")
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
	t.templ.Execute(w, r)
}

func main() {
	flag.Parse()
	r := chatroom.NewRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.Run()
	log.Println("Starting the server and listening on port", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("Server could not startup: %v", err)
	}
}
