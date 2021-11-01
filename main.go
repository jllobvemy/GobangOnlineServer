package main

import (
	"flag"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var addr = flag.String("addr", ":8082", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
// main
func main() {
	// ..........
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/num", func(writer http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("roomid")
		writer.Write([]byte(strconv.Itoa(len(hub.clients[id]))))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
