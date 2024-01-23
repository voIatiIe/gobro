package gobro

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type ServerOpts struct {
	addr string
	url string
}

type Server struct { Opts *ServerOpts }

type ServerOptsFunc func(*ServerOpts)

func WithAddr(addr string) ServerOptsFunc {
	return func(opts *ServerOpts) { opts.addr = addr }
}

func WithUrl(url string) ServerOptsFunc {
	return func(opts *ServerOpts) { opts.url = url }
}

func defaultServerOpts() *ServerOpts {
	return &ServerOpts{
		addr: ":8010",
		url: "/ws",
	}
}

func NewServer(opts ...ServerOptsFunc) *Server {
	defaultOpts := defaultServerOpts()

	// os.RemoveAll("./screenshots/")
	// os.MkdirAll("./screenshots/", 0o644)

	for _, opt := range opts { opt(defaultOpts) }
	
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	return &Server{Opts: defaultOpts}
}

func (s *Server) Start() error {
	log.Println("Starting server...")

	http.HandleFunc(s.Opts.url, WSHandler)

	return http.ListenAndServe(s.Opts.addr, nil)
}

type CursorMessage struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}


func WSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	browser, err := NewBrowser("https://google.com", WithQuality(50))
	if err != nil { 
		log.Println("Could not initialize browser:", err)
		return
	}
	defer browser.Cancel()

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Could not upgrade connection:", err)
		return
	}
	defer ws.Close()

	// ####

	// var message CursorMessage
	// for {
	// 	_, data, err := ws.ReadMessage()

	// 	if err != nil {
	// 		log.Println("Could not read message:", err)
	// 		return
	// 	}
	// 	if json.Unmarshal(data, &message) != nil {
	// 		log.Println("Could not parse message:", err)
	// 		return
	// 	}
	// 	log.Printf("Received: %+v\n", message)
	// }

	// ####

	ch := make(chan CursorMessage)
	wg := make(chan struct{}, 2)

	WGFunc := func(ch chan struct{}) { <-ch }

	go func() {
		defer log.Println("Exiting websocket pooling loop...")
		defer WGFunc(wg)
		var message CursorMessage

		for {
			select {
				case <-browser.Ctx.Done():
					return
				default:
					_, data, err := ws.ReadMessage()

					if err != nil {
						log.Println("Could not read message:", err)
						browser.Cancel()
						return
					}
					if json.Unmarshal(data, &message) != nil {
						log.Println("Could not parse message:", err)
						browser.Cancel()
						return
					}
					log.Printf("Received: %+v\n", message)

					ch <- message
			}
		}
	}()

	go func() {
		defer log.Println("Exiting browser rendering loop...")
		defer WGFunc(wg)

		for {
			select {
				case <-browser.Ctx.Done():
					return
				case message := <-ch:
					if err := browser.MoveCursor(message.X, message.Y); err != nil { browser.Cancel(); return }
					// if err := browser.Screenshot(); err != nil { browser.Cancel(); return }

					log.Println("[Render step]")
				// default:
				// 	if err := browser.Screenshot(); err != nil { browser.Cancel(); return }
			}
		}
	}()

	for range wg {}
}
