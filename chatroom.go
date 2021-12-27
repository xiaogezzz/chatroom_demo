package main

import (
	"chatroom/conn"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	http.HandleFunc("/chat", conn.Connection)
	http.HandleFunc("/join", join)
	http.HandleFunc("/", index)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func join(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	if uid == "" {
		t, _ := template.ParseFiles("./template/index.tpl")
		_ = t.Execute(w, nil)
		return
	}

	t, _ := template.ParseFiles("./template/room.tpl")
	log.Println(uid + " join")
	_ = t.Execute(w, map[string]string{"uid": uid})
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/index.tpl")
	_ = t.Execute(w, nil)
}
