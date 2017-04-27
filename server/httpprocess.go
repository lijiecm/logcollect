package server

import (
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HTTPStart() {
	http.HandleFunc("/log/ws", handleConnections)
	http.HandleFunc("/log/s", handleSelectArg)
	http.ListenAndServe(":8010", nil)
}

func handleSelectArg(w http.ResponseWriter, r *http.Request) {

}

// 收到要查日志文件，开始write
// switch log则重新write
func wsReader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	for {
		// 定义标记，标记变化则退出write goroutine
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func wsWriter(ws *websocket.Conn, filepath string) {
	defer ws.Close()

	for {
		t, err := tail.TailFile(filepath, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
		if err != nil {
			log.Printf("tail file failed, err: %v", err)
			break
		}
		log.Println("test")
		for line := range t.Lines {
			log.Println(line.Text)
			if err := ws.WriteMessage(websocket.TextMessage, []byte(line.Text)); err != nil {
				return
			}
		}
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Action,Module")
	}
	if r.Method == "OPTIONS" {
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("test")
		log.Fatal(err)
	}
	wsReader(ws)
}
