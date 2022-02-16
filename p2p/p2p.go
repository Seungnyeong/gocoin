package p2p

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snkim/sncoin/utils"
)


var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {

	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	// Port :3000 will upgrade the request from :4000
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	fmt.Printf("%s want to upgrade to port \n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	peer := initPeer(conn, ip, openPort)
	time.Sleep(10 * time.Second)
	peer.inbox <- []byte("Hello From 4000!")
}

func AddPeer(address, port, openPort string) {
	// from :4000 -> :3000
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}