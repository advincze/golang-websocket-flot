package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

var port *int = flag.Int("p", 23456, "Port to listen.")
var tick *time.Duration = flag.Duration("tick", time.Second, "Tick")
var openBrowser *bool = flag.Bool("open", false, "open browser")
var verbose *bool = flag.Bool("v", false, "verbose")

type hub struct {
	connections map[*websocket.Conn]bool
	name        string
}

var hubs map[string]*hub

func (h hub) register(conn *websocket.Conn) {
	h.connections[conn] = true
}

func (h hub) unregister(conn *websocket.Conn) {
	delete(h.connections, conn)
}

func (h hub) broadcast(message string) {
	var err error
	for ws := range h.connections {
		err = websocket.Message.Send(ws, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func wsServer(ws *websocket.Conn) {
	hub := hubs["time"] //TODO extract from ws config
	hub.register(ws)
	if *verbose {
		log.Println("registered for ws")
	}

	for {
		var buf string
		err := websocket.Message.Receive(ws, &buf)
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	hub.unregister(ws)
	if *verbose {
		log.Println("unregistered for ws")
	}
}

func tickAndDo(fn func(time.Time), d time.Duration) {
	ticker := time.NewTicker(d)
	for now := range ticker.C {
		fn(now)
		if *verbose {
			log.Println("tick", now.Format("15:04:05"))
		}
	}
}

type Event struct {
	TimeStamp   int64
	IntValue    int64
	StringValue string
}

func main() {
	flag.Parse()
	hubs = map[string]*hub{
		"time": &hub{connections: make(map[*websocket.Conn]bool)},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	go tickAndDo(func(now time.Time) {
		for ws := range hubs["time"].connections {
			event := &Event{TimeStamp: now.Unix(), IntValue: r.Int63n(200)}
			websocket.JSON.Send(ws, *event)
		}
	}, *tick)

	http.Handle("/ws", websocket.Handler(wsServer))
	http.Handle("/", http.FileServer(http.Dir("static")))
	if *openBrowser {
		openUrlInBrowser("http://localhost:" + strconv.Itoa(*port) + "/index.html")
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic("ListenANdServe: " + err.Error())
	}

}

func openUrlInBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command(`C:\Windows\System32\rundll32.exe`, "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Print("unsupported platform")
	}
	if err != nil {
		panic(err)
	}
}
