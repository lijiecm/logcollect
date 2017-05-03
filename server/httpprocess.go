package server

import (
	"github.com/gorilla/websocket"
	//"github.com/hpcloud/tail"
	"github.com/weihualiu/logcollect/store"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

	staticPage()

	http.ListenAndServe(":8010", nil)
}

func staticPage() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(filepath.Join(wd, "public/"))
	http.Handle("/log/", http.StripPrefix("/log/", http.FileServer(http.Dir(filepath.Join(wd, "public/")))))
}

func handleSelectArg(w http.ResponseWriter, r *http.Request) {
	//return text json format
	//  list data, flag next exist
	// 单个Tag列表获取, tag number
}

// 收到要查日志文件，开始write
// switch log则重新write
func wsReader(ws *websocket.Conn) {
	defer func() {
		ws.Close()
		log.Println("ws is closed!")
	}()
	ws.SetReadLimit(512)
	for {
		// 定义标记，标记变化则退出write goroutine
		_, content, err := ws.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("read.... %s", string(content))

		dataCh := make(chan []byte, 10)
		store.Monitors.Set(string(content), dataCh)

		currwd, err := os.Getwd()
		fpath := filepath.Join(currwd, "data/api/", string(content))
		go wsWriter(ws, fpath, dataCh)
	}
}

func wsWriter(ws *websocket.Conn, filepath string, dataCh chan []byte) {
	defer func() {
		ws.Close()
		close(dataCh)
		log.Println("ws is closed!")
	}()
	for {
		//t, err := tail.TailFile(filepath, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
		//if err != nil {
		//	log.Printf("tail file failed, err: %v", err)
		//	break
		//}
		//for line := range t.Lines {
		//	if err := ws.WriteMessage(websocket.TextMessage, []byte(line.Text)); err != nil {
		//		log.Println("read tail file end,", err.Error())
		//		return
		//	}
		//}
		time.Sleep(time.Millisecond * 100)
		select {
		case content := <-dataCh:
			if err := ws.WriteMessage(websocket.TextMessage, content); err != nil {
				log.Println("read tail file end,", err.Error())
				return
			}
		default:
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
