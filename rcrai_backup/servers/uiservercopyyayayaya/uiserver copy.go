package uiservercopyyayayaya

import (
	"dittor/model"
	"dittor/utils/refs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

// ----------------------------------------------------------------
// adapter!!
type Conf struct {
	Token          *string
	ProcessRunSign *int
	SourceId       *string
	Port           int
}

var conf = &Conf{
	Token:          nil,                                                           // ?
	ProcessRunSign: refs.Intref(0),                                                // ?
	SourceId:       refs.StringRef(fmt.Sprintf(`%x`, string(bson.NewObjectId()))), // ? TODO replace this.
	Port:           8990,
	// Port:           8789,
}

// ----------------------------------------------------------------

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ping(ws *websocket.Conn, wg *sync.WaitGroup) {
	defer (*wg).Done()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
				return
			}
		}
	}
}

func readMessage(ws *websocket.Conn) {
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		var read model.ElectronRead
		err = json.Unmarshal(message, &read)
		if err != nil {
			log.Println("read:", err)
			continue
		}
		if read.Action == "start" {
			conf.Token = &read.Token
		}
	}
}

func StartAndStopMessage(ws *websocket.Conn, wg *sync.WaitGroup) {
	defer (*wg).Done()
	defaultInt := 0
	for {
		if *conf.ProcessRunSign != defaultInt {
			defaultInt = *conf.ProcessRunSign
			if defaultInt == 1 {
				start, err := json.Marshal(&model.ElectronSend{
					Action:   "source_id",
					SourceId: *conf.SourceId,
				})
				if err != nil {
					continue
				}
				err = ws.WriteMessage(websocket.TextMessage, start)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				over, err := json.Marshal(&model.ElectronSend{
					Action:   "stop",
					SourceId: *conf.SourceId,
				})
				if err != nil {
					continue
				}
				err = ws.WriteMessage(websocket.TextMessage, over)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go Ping(c, &wg)
	go StartAndStopMessage(c, &wg)
	readMessage(c)
	wg.Wait()
}

func NewSocket() {
	http.HandleFunc("/v1/app/ws", serveWs)

	var addr = fmt.Sprintf("127.0.0.1:%d", conf.Port)
	fmt.Println(addr)

	err := http.ListenAndServe(addr, nil)
	fmt.Println("http exit!")
	if err != nil {
		log.Fatal(err)
	}
}
