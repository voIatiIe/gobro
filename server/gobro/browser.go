package gobro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/gorilla/websocket"
)

type BrowserOpts struct {
	width, height int64
	quality int
}

func defaultBrowserOpts() *BrowserOpts {
	return &BrowserOpts{
		width: 1920,
		height: 1080,
		quality: 100,
	}
}

type BrowserOptsFunc func(*BrowserOpts)

func WithQuality(quality int) BrowserOptsFunc { return func(opts *BrowserOpts) { opts.quality = quality } }

func WithWidth(width int64) BrowserOptsFunc { return func(opts *BrowserOpts) { opts.width = width } }

func WithHeight(height int64) BrowserOptsFunc { return func(opts *BrowserOpts) { opts.height = height } }

type Browser struct { 
	Opts *BrowserOpts
	Ctx context.Context
	Cancel context.CancelFunc
}

func NewBrowser(url string, opts ...BrowserOptsFunc) (*Browser, error) {
	defaultOpts := defaultBrowserOpts()

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Headless,
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), allocOpts...)

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	for _, opt := range opts { opt(defaultOpts) }

	browser := &Browser{
		Opts: defaultOpts,
		Ctx: ctx,
		Cancel: cancel,
	}

	tasks := chromedp.Tasks{
		chromedp.EmulateViewport(browser.Opts.width, browser.Opts.height),
		chromedp.Navigate(url),
		// InitCursor(),
	}

	if err := chromedp.Run(browser.Ctx, tasks); err != nil { return nil, err }

	return browser, nil
}

func (b *Browser) GenericClickFunc(x, y float64, opts ...chromedp.MouseOption) error {
	x *= float64(b.Opts.width)
	y *= float64(b.Opts.height)

	action := MouseClickXY(x, y, opts...)

	return chromedp.Run(b.Ctx, action)
}

func (b *Browser) MoveFunc(x, y float64) error { return b.GenericClickFunc(x, y, chromedp.ButtonNone) }
func (b *Browser) LeftClickFunc(x, y float64) error { return b.GenericClickFunc(x, y, chromedp.ButtonLeft) }

func (b *Browser) TakeScreenshot(buffer *[]byte, quality int) error {
	action := Screenshot(buffer, quality)

	return chromedp.Run(b.Ctx, action)
}

func (b *Browser) Stream(ws *websocket.Conn, wg *sync.WaitGroup, lock sync.Locker) {
	defer wg.Done()
	defer log.Println("Terminating stream...")

	var buffer []byte

	for {
		select {
			case <-b.Ctx.Done():
				return
			default:
				lock.Lock()
				if err := b.TakeScreenshot(&buffer, b.Opts.quality); err != nil {
					log.Println("Could not take screenshot:", err)
					b.Cancel()
					lock.Unlock()
					return
				}
				lock.Unlock()

				if err := ws.WriteMessage(websocket.BinaryMessage, buffer); err != nil {
					log.Println("Could not write message:", err)
					b.Cancel()
					return
				}
		}
	}
}

type CommandType string

const (
	Move CommandType = "MOUSE_MOVE"
	LeftClick CommandType = "MOUSE_CLICK"
	Scroll CommandType = "SCROLL"
	Resize CommandType = "RESIZE"

	Input CommandType = "INPUT"
	Delete CommandType = "DELETE"
)

type CommandPayload struct {
	Coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"coordinates"`
}

type CommandMessage struct {
	Type CommandType `json:"type"`
	Payload CommandPayload `json:"payload"`
}


func (b *Browser) Control(ws *websocket.Conn, wg *sync.WaitGroup, lock sync.Locker) {
	defer wg.Done()

	var message CommandMessage

	for {
		select {
			case <-b.Ctx.Done():
				return
			default:
				_, data, err := ws.ReadMessage()

				if err != nil {
					log.Println("Could not read message:", err)
					b.Cancel()
					return
				}

				if err := json.Unmarshal(data, &message); err != nil {
					log.Println("Could not unmarshal message:", err)
					b.Cancel()
					return
				}
				
				lock.Lock()
				if err := b.Execute(message); err != nil {
					log.Println("Could not execute command, skipped:", err)
				}
				lock.Unlock()
		}
	}
}

func (b *Browser) Execute(cmd CommandMessage) error {
	switch cmd.Type {
	case Move:
		return b.MoveFunc(cmd.Payload.Coordinates.X, cmd.Payload.Coordinates.Y)
	case LeftClick:
		return b.LeftClickFunc(cmd.Payload.Coordinates.X, cmd.Payload.Coordinates.Y)
	default:
		errStr := fmt.Sprintf("Unsupported command type: %s\n", cmd.Type)
		return errors.New(errStr)
	}
}
