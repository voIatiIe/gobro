package gobro

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ServerOpts struct {
	addr string
	url  string
}

type Server struct{ Opts *ServerOpts }

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
		url:  "/ws",
	}
}

func NewServer(opts ...ServerOptsFunc) *Server {
	defaultOpts := defaultServerOpts()

	for _, opt := range opts {
		opt(defaultOpts)
	}

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

	browser, err := NewBrowser(
		"https://google.com/",
		WithQuality(80),
		WithHeight(720),
		WithWidth(1280),
	)
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

	wg := &sync.WaitGroup{}
	wg.Add(2)
	lock := &sync.Mutex{}

	go browser.Control(ws, wg, lock)
	go browser.Stream(ws, wg, lock)

	wg.Wait()
}
